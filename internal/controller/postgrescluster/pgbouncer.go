/*
 Copyright 2021 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package postgrescluster

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crunchydata/postgres-operator/internal/logging"
	"github.com/crunchydata/postgres-operator/internal/naming"
	"github.com/crunchydata/postgres-operator/internal/pgbouncer"
	"github.com/crunchydata/postgres-operator/internal/postgres"
	"github.com/crunchydata/postgres-operator/pkg/apis/postgres-operator.crunchydata.com/v1alpha1"
)

// reconcilePGBouncer writes the objects necessary to run a PgBouncer Pod.
func (r *Reconciler) reconcilePGBouncer(
	ctx context.Context, cluster *v1alpha1.PostgresCluster,
) error {
	var (
		configmap *corev1.ConfigMap
		secret    *corev1.Secret
	)

	err := r.reconcilePGBouncerService(ctx, cluster)
	if err == nil {
		configmap, err = r.reconcilePGBouncerConfigMap(ctx, cluster)
	}
	if err == nil {
		secret, err = r.reconcilePGBouncerSecret(ctx, cluster)
	}
	if err == nil {
		err = r.reconcilePGBouncerDeployment(ctx, cluster, configmap, secret)
	}
	if err == nil {
		err = r.reconcilePGBouncerInPostgreSQL(ctx, cluster, secret)
	}
	return err
}

// +kubebuilder:rbac:resources=configmaps,verbs=get;delete;patch

// reconcilePGBouncerConfigMap writes the ConfigMap for a PgBouncer Pod.
func (r *Reconciler) reconcilePGBouncerConfigMap(
	ctx context.Context, cluster *v1alpha1.PostgresCluster,
) (*corev1.ConfigMap, error) {
	configmap := &corev1.ConfigMap{ObjectMeta: naming.ClusterPGBouncer(cluster)}
	configmap.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	if cluster.Spec.Proxy == nil || cluster.Spec.Proxy.PGBouncer == nil {
		// PgBouncer is disabled; delete the ConfigMap if it exists. Check the
		// client cache first using Get.
		key := client.ObjectKeyFromObject(configmap)
		err := errors.WithStack(r.Client.Get(ctx, key, &corev1.ConfigMap{}))
		if err == nil {
			err = errors.WithStack(r.Client.Delete(ctx, configmap))
		}
		return nil, client.IgnoreNotFound(err)
	}

	err := errors.WithStack(r.setControllerReference(cluster, configmap))

	configmap.Labels = map[string]string{
		naming.LabelCluster: cluster.Name,
		naming.LabelRole:    naming.RolePGBouncer,
	}

	if err == nil {
		pgbouncer.ConfigMap(cluster, configmap)
	}
	if err == nil {
		err = errors.WithStack(r.apply(ctx, configmap))
	}

	return configmap, err
}

// +kubebuilder:rbac:resources=pods,verbs=get;list

// reconcilePGBouncerInPostgreSQL writes the user and other objects needed by
// PgBouncer inside of PostgreSQL.
func (r *Reconciler) reconcilePGBouncerInPostgreSQL(
	ctx context.Context, cluster *v1alpha1.PostgresCluster, clusterSecret *corev1.Secret,
) error {
	if cluster.Status.Patroni == nil || cluster.Status.Patroni.SystemIdentifier == "" {
		// Patroni has not yet bootstrapped; there is nothing to do.
		// NOTE(cbandy): "Patroni bootstrapped" may not be the right check here.
		// The following code needs to be able to execute SQL that writes
		// objects in every PostgreSQL database (probably as the superuser).
		return nil
	}

	// Patroni has bootstrapped. Prepare to either add or remove PgBouncer from
	// PostgreSQL.

	action := func(ctx context.Context, exec postgres.Executor) error {
		return errors.WithStack(pgbouncer.EnableInPostgreSQL(ctx, exec, clusterSecret))
	}
	if cluster.Spec.Proxy == nil || cluster.Spec.Proxy.PGBouncer == nil {
		// PgBouncer is disabled.
		action = func(ctx context.Context, exec postgres.Executor) error {
			return errors.WithStack(pgbouncer.DisableInPostgreSQL(ctx, exec))
		}
	}

	// First, calculate a hash of the SQL that should be executed in PostgreSQL.

	revision, err := safeHash32(func(hasher io.Writer) error {
		// Discard log messages from the pgbouncer package about executing SQL.
		// Nothing is being "executed" yet.
		return action(logging.NewContext(ctx, logging.Discard()), func(
			_ context.Context, stdin io.Reader, _, _ io.Writer, command ...string,
		) error {
			_, err := io.Copy(hasher, stdin)
			if err == nil {
				_, err = fmt.Fprint(hasher, command)
			}
			return err
		})
	})
	if err != nil {
		return err
	}

	if revision == cluster.Status.Proxy.PGBouncer.PostgreSQLRevision {
		// The necessary SQL has already been applied; there's nothing more to do.

		// TODO(cbandy): Give the user a way to trigger execution regardless.
		// The value of an annotation could influence the hash, for example.
		return nil
	}

	// The necessary SQL has not been applied. Find a pod that can write to cluster.

	pods := &corev1.PodList{}
	instances, err := naming.AsSelector(naming.ClusterInstances(cluster.Name))
	if err == nil {
		err = errors.WithStack(
			r.Client.List(ctx, pods,
				client.InNamespace(cluster.Namespace),
				client.MatchingLabelsSelector{Selector: instances},
			))
	}

	var pod *corev1.Pod
	if err == nil {
		for i := range pods.Items {
			if pods.Items[i].Labels[naming.LabelRole] == naming.RolePatroniLeader {
				pod = &pods.Items[i]
				break
			}
		}
		if pod == nil {
			err = errors.New("could not find primary pod")
		}
	}

	// Apply the necessary SQL and record its hash in cluster.Status. Include
	// the hash in any log messages.

	if err == nil {
		ctx := logging.NewContext(ctx, logging.FromContext(ctx).WithValues("revision", revision))
		err = action(ctx, func(_ context.Context, stdin io.Reader, stdout, stderr io.Writer, command ...string) error {
			return r.PodExec(pod.Namespace, pod.Name, naming.ContainerDatabase, stdin, stdout, stderr, command...)
		})
	}
	if err == nil {
		cluster.Status.Proxy.PGBouncer.PostgreSQLRevision = revision
	}

	return err
}

// +kubebuilder:rbac:resources=secrets,verbs=get;delete;patch

// reconcilePGBouncerSecret writes the Secret for a PgBouncer Pod.
func (r *Reconciler) reconcilePGBouncerSecret(
	ctx context.Context, cluster *v1alpha1.PostgresCluster,
) (*corev1.Secret, error) {
	existing := &corev1.Secret{ObjectMeta: naming.ClusterPGBouncer(cluster)}
	err := errors.WithStack(
		r.Client.Get(ctx, client.ObjectKeyFromObject(existing), existing))
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	if cluster.Spec.Proxy == nil || cluster.Spec.Proxy.PGBouncer == nil {
		// PgBouncer is disabled; delete the Secret if it exists.
		if err == nil {
			err = errors.WithStack(r.Client.Delete(ctx, existing))
		}
		return nil, client.IgnoreNotFound(err)
	}

	err = client.IgnoreNotFound(err)

	intent := &corev1.Secret{ObjectMeta: naming.ClusterPGBouncer(cluster)}
	intent.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	intent.Type = corev1.SecretTypeOpaque

	if err == nil {
		err = errors.WithStack(r.setControllerReference(cluster, intent))
	}

	intent.Labels = map[string]string{
		naming.LabelCluster: cluster.Name,
		naming.LabelRole:    naming.RolePGBouncer,
	}

	if err == nil {
		err = pgbouncer.Secret(existing, intent)
	}
	if err == nil {
		err = errors.WithStack(r.apply(ctx, intent))
	}

	return intent, err
}

// +kubebuilder:rbac:resources=services,verbs=get;delete;patch

// reconcilePGBouncerService writes the Service that resolves to PgBouncer.
func (r *Reconciler) reconcilePGBouncerService(
	ctx context.Context, cluster *v1alpha1.PostgresCluster,
) error {
	service := &corev1.Service{ObjectMeta: naming.ClusterPGBouncer(cluster)}
	service.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Service"))

	if cluster.Spec.Proxy == nil || cluster.Spec.Proxy.PGBouncer == nil {
		// PgBouncer is disabled; delete the Service if it exists. Check the client
		// cache first using Get.
		key := client.ObjectKeyFromObject(service)
		err := errors.WithStack(r.Client.Get(ctx, key, &corev1.Service{}))
		if err == nil {
			err = errors.WithStack(r.Client.Delete(ctx, service))
		}
		return client.IgnoreNotFound(err)
	}

	err := errors.WithStack(r.setControllerReference(cluster, service))

	service.Labels = map[string]string{
		naming.LabelCluster: cluster.Name,
		naming.LabelRole:    naming.RolePGBouncer,
	}

	// Allocate an IP address and let Kubernetes manage the Endpoints by selecting
	// Pods with the PgBouncer role.
	// - https://docs.k8s.io/concepts/services-networking/service/#defining-a-service
	service.Spec.Type = corev1.ServiceTypeClusterIP
	service.Spec.Selector = map[string]string{
		naming.LabelCluster: cluster.Name,
		naming.LabelRole:    naming.RolePGBouncer,
	}

	// The TargetPort must be the name (not the number) of the PgBouncer
	// ContainerPort. This name allows the port number to differ between Pods,
	// which can happen during a rolling update.
	service.Spec.Ports = []corev1.ServicePort{{
		Name:       naming.PortPGBouncer,
		Port:       *cluster.Spec.Proxy.PGBouncer.Port,
		Protocol:   corev1.ProtocolTCP,
		TargetPort: intstr.FromString(naming.PortPGBouncer),
	}}

	if err == nil {
		err = errors.WithStack(r.apply(ctx, service))
	}

	return err
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;delete;patch

// reconcilePGBouncerDeployment writes the Deployment that runs PgBouncer.
func (r *Reconciler) reconcilePGBouncerDeployment(
	ctx context.Context, cluster *v1alpha1.PostgresCluster,
	configmap *corev1.ConfigMap, secret *corev1.Secret,
) error {
	deploy := &appsv1.Deployment{ObjectMeta: naming.ClusterPGBouncer(cluster)}
	deploy.SetGroupVersionKind(appsv1.SchemeGroupVersion.WithKind("Deployment"))

	// Set observations whether the deployment exists or not.
	defer func() {
		cluster.Status.Proxy.PGBouncer.Replicas = deploy.Status.Replicas
		cluster.Status.Proxy.PGBouncer.ReadyReplicas = deploy.Status.ReadyReplicas

		// NOTE(cbandy): This should be somewhere else when there is more than
		// one proxy implementation.

		var available *appsv1.DeploymentCondition
		for i := range deploy.Status.Conditions {
			if deploy.Status.Conditions[i].Type == appsv1.DeploymentAvailable {
				available = &deploy.Status.Conditions[i]
			}
		}

		if available == nil {
			// Avoid a panic! Fixed in Kubernetes v1.21.0 and controller-runtime v0.9.0-alpha.0.
			// - https://issue.k8s.io/99714
			if len(cluster.Status.Conditions) > 0 {
				meta.RemoveStatusCondition(&cluster.Status.Conditions, v1alpha1.ProxyAvailable)
			}
		} else {
			meta.SetStatusCondition(&cluster.Status.Conditions, metav1.Condition{
				Type:    v1alpha1.ProxyAvailable,
				Status:  metav1.ConditionStatus(available.Status),
				Reason:  available.Reason,
				Message: available.Message,

				LastTransitionTime: available.LastTransitionTime,
				ObservedGeneration: cluster.Generation,
			})
		}
	}()

	if cluster.Spec.Proxy == nil || cluster.Spec.Proxy.PGBouncer == nil {
		// PgBouncer is disabled; delete the Deployment if it exists. Check the
		// client cache first using Get.
		key := client.ObjectKeyFromObject(deploy)
		err := errors.WithStack(r.Client.Get(ctx, key, deploy))
		if err == nil {
			err = errors.WithStack(r.Client.Delete(ctx, deploy))
		}
		return client.IgnoreNotFound(err)
	}

	err := errors.WithStack(r.setControllerReference(cluster, deploy))

	deploy.Labels = map[string]string{
		naming.LabelCluster: cluster.Name,
		naming.LabelRole:    naming.RolePGBouncer,
	}
	deploy.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			naming.LabelCluster: cluster.Name,
			naming.LabelRole:    naming.RolePGBouncer,
		},
	}
	deploy.Spec.Template.Labels = map[string]string{
		naming.LabelCluster: cluster.Name,
		naming.LabelRole:    naming.RolePGBouncer,
	}

	deploy.Spec.Replicas = new(int32)
	*deploy.Spec.Replicas = 1

	// Don't clutter the namespace with extra ReplicaSets.
	deploy.Spec.RevisionHistoryLimit = new(int32) // zero

	// TODO(cbandy): Consider the desired rollout behavior. The defaults here
	// with "spec.replicas=1" cause a surge of one before removing the one old
	// pod.
	// - https://docs.k8s.io/concepts/workloads/controllers/deployment/#rolling-update-deployment
	// - https://docs.k8s.io/concepts/workloads/controllers/statefulset/#on-delete
	deploy.Spec.Strategy.Type = appsv1.RollingUpdateDeploymentStrategyType
	deploy.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{}

	// Use scheduling constraints from the cluster spec.
	deploy.Spec.Template.Spec.Affinity = cluster.Spec.Proxy.PGBouncer.Affinity
	deploy.Spec.Template.Spec.Tolerations = cluster.Spec.Proxy.PGBouncer.Tolerations

	// Restart containers any time they stop, die, are killed, etc.
	// - https://docs.k8s.io/concepts/workloads/pods/pod-lifecycle/#restart-policy
	deploy.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways

	// ShareProcessNamespace makes Kubernetes' pause process PID 1 and lets
	// containers see each other's processes.
	// - https://docs.k8s.io/tasks/configure-pod-container/share-process-namespace/
	deploy.Spec.Template.Spec.ShareProcessNamespace = new(bool)
	*deploy.Spec.Template.Spec.ShareProcessNamespace = true

	// There's no need for individual DNS names of PgBouncer pods.
	deploy.Spec.Template.Spec.Subdomain = ""

	// PgBouncer does not make any Kubernetes API calls. Use the default
	// ServiceAccount and do not mount its credentials.
	deploy.Spec.Template.Spec.AutomountServiceAccountToken = new(bool) // false

	True := true
	deploy.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot: &True,
	}

	if err == nil {
		pgbouncer.Pod(cluster, configmap, secret, &deploy.Spec.Template.Spec)
	}
	if err == nil {
		err = errors.WithStack(r.apply(ctx, deploy))
	}

	return err
}
