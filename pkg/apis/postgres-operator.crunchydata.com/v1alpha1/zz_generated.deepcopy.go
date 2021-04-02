// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Archive) DeepCopyInto(out *Archive) {
	*out = *in
	in.PGBackRest.DeepCopyInto(&out.PGBackRest)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Archive.
func (in *Archive) DeepCopy() *Archive {
	if in == nil {
		return nil
	}
	out := new(Archive)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DedicatedRepo) DeepCopyInto(out *DedicatedRepo) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DedicatedRepo.
func (in *DedicatedRepo) DeepCopy() *DedicatedRepo {
	if in == nil {
		return nil
	}
	out := new(DedicatedRepo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PGBackRestArchive) DeepCopyInto(out *PGBackRestArchive) {
	*out = *in
	if in.Configuration != nil {
		in, out := &in.Configuration, &out.Configuration
		*out = make([]v1.VolumeProjection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Repos != nil {
		in, out := &in.Repos, &out.Repos
		*out = make([]RepoVolume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RepoHost != nil {
		in, out := &in.RepoHost, &out.RepoHost
		*out = new(RepoHost)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PGBackRestArchive.
func (in *PGBackRestArchive) DeepCopy() *PGBackRestArchive {
	if in == nil {
		return nil
	}
	out := new(PGBackRestArchive)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PGBackRestStatus) DeepCopyInto(out *PGBackRestStatus) {
	*out = *in
	if in.RepoHost != nil {
		in, out := &in.RepoHost, &out.RepoHost
		*out = new(RepoHostStatus)
		**out = **in
	}
	if in.Repos != nil {
		in, out := &in.Repos, &out.Repos
		*out = make([]RepoVolumeStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PGBackRestStatus.
func (in *PGBackRestStatus) DeepCopy() *PGBackRestStatus {
	if in == nil {
		return nil
	}
	out := new(PGBackRestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PGBouncerPodSpec) DeepCopyInto(out *PGBouncerPodSpec) {
	*out = *in
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PGBouncerPodSpec.
func (in *PGBouncerPodSpec) DeepCopy() *PGBouncerPodSpec {
	if in == nil {
		return nil
	}
	out := new(PGBouncerPodSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PGBouncerPodStatus) DeepCopyInto(out *PGBouncerPodStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PGBouncerPodStatus.
func (in *PGBouncerPodStatus) DeepCopy() *PGBouncerPodStatus {
	if in == nil {
		return nil
	}
	out := new(PGBouncerPodStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PatroniSpec) DeepCopyInto(out *PatroniSpec) {
	*out = *in
	in.DynamicConfiguration.DeepCopyInto(&out.DynamicConfiguration)
	if in.LeaderLeaseDurationSeconds != nil {
		in, out := &in.LeaderLeaseDurationSeconds, &out.LeaderLeaseDurationSeconds
		*out = new(int32)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.SyncPeriodSeconds != nil {
		in, out := &in.SyncPeriodSeconds, &out.SyncPeriodSeconds
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PatroniSpec.
func (in *PatroniSpec) DeepCopy() *PatroniSpec {
	if in == nil {
		return nil
	}
	out := new(PatroniSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PatroniStatus) DeepCopyInto(out *PatroniStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PatroniStatus.
func (in *PatroniStatus) DeepCopy() *PatroniStatus {
	if in == nil {
		return nil
	}
	out := new(PatroniStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresCluster) DeepCopyInto(out *PostgresCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresCluster.
func (in *PostgresCluster) DeepCopy() *PostgresCluster {
	if in == nil {
		return nil
	}
	out := new(PostgresCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PostgresCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresClusterList) DeepCopyInto(out *PostgresClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PostgresCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresClusterList.
func (in *PostgresClusterList) DeepCopy() *PostgresClusterList {
	if in == nil {
		return nil
	}
	out := new(PostgresClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PostgresClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresClusterSpec) DeepCopyInto(out *PostgresClusterSpec) {
	*out = *in
	in.Archive.DeepCopyInto(&out.Archive)
	if in.CustomTLSSecret != nil {
		in, out := &in.CustomTLSSecret, &out.CustomTLSSecret
		*out = new(v1.SecretProjection)
		(*in).DeepCopyInto(*out)
	}
	if in.InstanceSets != nil {
		in, out := &in.InstanceSets, &out.InstanceSets
		*out = make([]PostgresInstanceSetSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.OpenShift != nil {
		in, out := &in.OpenShift, &out.OpenShift
		*out = new(bool)
		**out = **in
	}
	if in.Patroni != nil {
		in, out := &in.Patroni, &out.Patroni
		*out = new(PatroniSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.Proxy != nil {
		in, out := &in.Proxy, &out.Proxy
		*out = new(PostgresProxySpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresClusterSpec.
func (in *PostgresClusterSpec) DeepCopy() *PostgresClusterSpec {
	if in == nil {
		return nil
	}
	out := new(PostgresClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresClusterStatus) DeepCopyInto(out *PostgresClusterStatus) {
	*out = *in
	if in.Patroni != nil {
		in, out := &in.Patroni, &out.Patroni
		*out = new(PatroniStatus)
		**out = **in
	}
	if in.PGBackRest != nil {
		in, out := &in.PGBackRest, &out.PGBackRest
		*out = new(PGBackRestStatus)
		(*in).DeepCopyInto(*out)
	}
	out.Proxy = in.Proxy
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresClusterStatus.
func (in *PostgresClusterStatus) DeepCopy() *PostgresClusterStatus {
	if in == nil {
		return nil
	}
	out := new(PostgresClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresInstanceSetSpec) DeepCopyInto(out *PostgresInstanceSetSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.Resources.DeepCopyInto(&out.Resources)
	in.VolumeClaimSpec.DeepCopyInto(&out.VolumeClaimSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresInstanceSetSpec.
func (in *PostgresInstanceSetSpec) DeepCopy() *PostgresInstanceSetSpec {
	if in == nil {
		return nil
	}
	out := new(PostgresInstanceSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresProxySpec) DeepCopyInto(out *PostgresProxySpec) {
	*out = *in
	if in.PGBouncer != nil {
		in, out := &in.PGBouncer, &out.PGBouncer
		*out = new(PGBouncerPodSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresProxySpec.
func (in *PostgresProxySpec) DeepCopy() *PostgresProxySpec {
	if in == nil {
		return nil
	}
	out := new(PostgresProxySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresProxyStatus) DeepCopyInto(out *PostgresProxyStatus) {
	*out = *in
	out.PGBouncer = in.PGBouncer
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresProxyStatus.
func (in *PostgresProxyStatus) DeepCopy() *PostgresProxyStatus {
	if in == nil {
		return nil
	}
	out := new(PostgresProxyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepoHost) DeepCopyInto(out *RepoHost) {
	*out = *in
	if in.Dedicated != nil {
		in, out := &in.Dedicated, &out.Dedicated
		*out = new(DedicatedRepo)
		(*in).DeepCopyInto(*out)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.SSHConfiguration != nil {
		in, out := &in.SSHConfiguration, &out.SSHConfiguration
		*out = new(v1.ConfigMapProjection)
		(*in).DeepCopyInto(*out)
	}
	if in.SSHSecret != nil {
		in, out := &in.SSHSecret, &out.SSHSecret
		*out = new(v1.SecretProjection)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepoHost.
func (in *RepoHost) DeepCopy() *RepoHost {
	if in == nil {
		return nil
	}
	out := new(RepoHost)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepoHostStatus) DeepCopyInto(out *RepoHostStatus) {
	*out = *in
	out.TypeMeta = in.TypeMeta
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepoHostStatus.
func (in *RepoHostStatus) DeepCopy() *RepoHostStatus {
	if in == nil {
		return nil
	}
	out := new(RepoHostStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepoVolume) DeepCopyInto(out *RepoVolume) {
	*out = *in
	in.VolumeClaimSpec.DeepCopyInto(&out.VolumeClaimSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepoVolume.
func (in *RepoVolume) DeepCopy() *RepoVolume {
	if in == nil {
		return nil
	}
	out := new(RepoVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepoVolumeStatus) DeepCopyInto(out *RepoVolumeStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepoVolumeStatus.
func (in *RepoVolumeStatus) DeepCopy() *RepoVolumeStatus {
	if in == nil {
		return nil
	}
	out := new(RepoVolumeStatus)
	in.DeepCopyInto(out)
	return out
}
