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

package postgres

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

// This example demonstrates how Executor can work with exec.Cmd.
func ExampleExecutor_execCmd() {
	_ = Executor(func(
		ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, command ...string,
	) error {
		// #nosec G204 Nothing calls the function defined in this example.
		cmd := exec.CommandContext(ctx, command[0], command[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = stdin, stdout, stderr
		return cmd.Run()
	})
}

func TestExecutorExecInDatabasesFromQuery(t *testing.T) {
	expected := errors.New("splat")
	fn := func(
		_ context.Context, stdin io.Reader, stdout, stderr io.Writer, command ...string,
	) error {
		b, err := ioutil.ReadAll(stdin)
		assert.NilError(t, err)
		assert.Equal(t, string(b), `statements; to run;`)

		assert.DeepEqual(t, command, []string{
			"bash", "-ceu", "--", `
sql_target=$(< /dev/stdin)
sql_databases="$1"
shift 1

databases=$(psql "$@" -X -Aqt --file=- <<< "${sql_databases}")
while read -r database; do
	psql "$@" -X --file=- "${database}" <<< "${sql_target}"
done <<< "${databases}"
`,
			"-",
			`db query`,
			"--set=CASE=sEnSiTiVe",
			"--set=different=vars",
			"--set=lots=of",
		})

		_, _ = io.WriteString(stdout, "some stdout")
		_, _ = io.WriteString(stderr, "and stderr")
		return expected
	}

	stdout, stderr, err := Executor(fn).ExecInDatabasesFromQuery(
		context.Background(), `db query`, `statements; to run;`, map[string]string{
			"lots":      "of",
			"different": "vars",
			"CASE":      "sEnSiTiVe",
		})

	assert.Equal(t, expected, err, "expected function to be called")
	assert.Equal(t, stdout, "some stdout")
	assert.Equal(t, stderr, "and stderr")

	t.Run("ShellCheck", func(t *testing.T) {
		shellcheck, err := exec.LookPath("shellcheck")
		if err != nil {
			t.Skip(`requires "shellcheck" executable`)
		} else {
			output, err := exec.Command(shellcheck, "--version").CombinedOutput()
			assert.NilError(t, err)
			t.Logf("using %q:\n%s", shellcheck, output)
		}

		_, _, _ = Executor(func(
			_ context.Context, _ io.Reader, _, _ io.Writer, command ...string,
		) error {
			// Expect a bash command with an inline script.
			assert.DeepEqual(t, command[:3], []string{"bash", "-ceu", "--"})
			assert.Assert(t, len(command) > 3)
			script := command[3]

			// Write out that inline script.
			dir := t.TempDir()
			file := filepath.Join(dir, "script.bash")
			assert.NilError(t, ioutil.WriteFile(file, []byte(script), 0o600))

			// Expect shellcheck to be happy.
			cmd := exec.Command(shellcheck, "--enable=all", file)
			output, err := cmd.CombinedOutput()
			assert.NilError(t, err, "%q\n%s", cmd.Args, output)

			return nil
		}).ExecInDatabasesFromQuery(context.Background(), "", "", nil)
	})
}
