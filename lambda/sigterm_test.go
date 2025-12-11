//go:build go1.15
// +build go1.15

package lambda

import (
	"fmt"
	"io"
	"io/ioutil" //nolint: staticcheck
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnableSigterm(t *testing.T) {
	containerCmd := ""
	if _, err := exec.LookPath("finch"); err == nil {
		containerCmd = "finch"
	} else if _, err := exec.LookPath("docker"); err == nil {
		containerCmd = "docker"
	} else {
		t.Skip("finch or docker required")
	}

	testDir := t.TempDir()

	// compile our handler, it'll always run to timeout ensuring the SIGTERM is triggered
	handlerBuild := exec.Command("go", "build", "-o", filepath.Join(testDir, "bootstrap"), "./testdata/sigterm.go")
	handlerBuild.Env = append(os.Environ(), "GOOS=linux")
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
			// Find an available port
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			port := listener.Addr().(*net.TCPAddr).Port
			listener.Close()

			cmdArgs := []string{"run", "--rm",
				"-v", testDir + ":/var/runtime:ro,delegated",
				"-p", fmt.Sprintf("%d:8080", port),
				"-e", "AWS_LAMBDA_FUNCTION_TIMEOUT=2"}
			for _, env := range opts.envVars {
				cmdArgs = append(cmdArgs, "-e", env)
			}
			cmdArgs = append(cmdArgs, "public.ecr.aws/lambda/provided:al2023", "bootstrap")

			cmd := exec.Command(containerCmd, cmdArgs...)
			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)
			stderr, err := cmd.StderrPipe()
			require.NoError(t, err)

			var logBuf strings.Builder
			logDone := make(chan struct{})
			go func() {
				_, _ = io.Copy(io.MultiWriter(os.Stderr, &logBuf), io.MultiReader(stdout, stderr))
				close(logDone)
			}()

			require.NoError(t, cmd.Start())
			t.Cleanup(func() { _ = cmd.Process.Kill() })

			time.Sleep(5 * time.Second) // Wait for container to start

			client := &http.Client{Timeout: 5 * time.Second}
			invokeURL := fmt.Sprintf("http://127.0.0.1:%d/2015-03-31/functions/function/invocations", port)
			resp, err := client.Post(invokeURL, "application/json", strings.NewReader("{}"))
			require.NoError(t, err)
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Task timed out after 2.00 seconds", string(body))

			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			<-logDone

			logs := logBuf.String()
			t.Logf("stdout:\n%s", logs)
			opts.assertLogs(t, logs)
		})
	}
}
