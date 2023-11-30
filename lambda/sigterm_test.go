//go:build go1.15
// +build go1.15

package lambda

import (
	"io/ioutil" //nolint: staticcheck
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnableSigterm(t *testing.T) {
	if _, err := exec.LookPath("aws-lambda-rie"); err != nil {
		t.Skipf("%v - install from https://github.com/aws/aws-lambda-runtime-interface-emulator/", err)
	}

	testDir := t.TempDir()

	// compile our handler, it'll always run to timeout ensuring the SIGTERM is triggered by aws-lambda-rie
	handlerBuild := exec.Command("go", "build", "-o", path.Join(testDir, "sigterm.handler"), "./testdata/sigterm.go")
	handlerBuild.Stderr = os.Stderr
	handlerBuild.Stdout = os.Stderr
	require.NoError(t, handlerBuild.Run())

	for name, opts := range map[string]struct {
		envVars    []string
		assertLogs func(t *testing.T, logs string)
	}{
		"baseline": {
			assertLogs: func(t *testing.T, logs string) {
				assert.NotContains(t, logs, "Hello SIGTERM!")
				assert.NotContains(t, logs, "I've been TERMINATED!")
			},
		},
		"sigterm enabled": {
			envVars: []string{"ENABLE_SIGTERM=please"},
			assertLogs: func(t *testing.T, logs string) {
				assert.Contains(t, logs, "Hello SIGTERM!")
				assert.Contains(t, logs, "I've been TERMINATED!")
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			addr1 := "localhost:" + strconv.Itoa(8000+rand.Intn(999))
			addr2 := "localhost:" + strconv.Itoa(9000+rand.Intn(999))
			rieInvokeAPI := "http://" + addr1 + "/2015-03-31/functions/function/invocations"
			// run the runtime interface emulator, capture the logs for assertion
			cmd := exec.Command("aws-lambda-rie", "--runtime-interface-emulator-address", addr1, "--runtime-api-address", addr2, "sigterm.handler")
			cmd.Env = append([]string{
				"PATH=" + testDir,
				"AWS_LAMBDA_FUNCTION_TIMEOUT=2",
			}, opts.envVars...)
			cmd.Stderr = os.Stderr
			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)
			var logs string
			done := make(chan interface{}) // closed on completion of log flush
			go func() {
				logBytes, err := ioutil.ReadAll(stdout)
				require.NoError(t, err)
				logs = string(logBytes)
				close(done)
			}()
			require.NoError(t, cmd.Start())
			t.Cleanup(func() { _ = cmd.Process.Kill() })

			// give a moment for the port to bind
			time.Sleep(500 * time.Millisecond)

			client := &http.Client{Timeout: 5 * time.Second} // http client timeout to prevent case from hanging on aws-lambda-rie
			resp, err := client.Post(rieInvokeAPI, "application/json", strings.NewReader("{}"))
			require.NoError(t, err)
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, string(body), "Task timed out after 2.00 seconds")

			require.NoError(t, cmd.Process.Kill()) // now ensure the logs are drained
			<-done
			t.Logf("stdout:\n%s", logs)
			opts.assertLogs(t, logs)
		})
	}
}
