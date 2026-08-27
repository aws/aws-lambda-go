//go:build go1.15
// +build go1.15

package lambda

import (
	"fmt"
	"io"
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

	// Pre-pull the container image so that pull latency doesn't count against
	// the per-subtest readiness deadline.
	pull := exec.Command(containerCmd, "pull", "public.ecr.aws/lambda/provided:al2023")
	pull.Stdout = os.Stderr
	pull.Stderr = os.Stderr
	require.NoError(t, pull.Run())

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
				"-e", "AWS_LAMBDA_FUNCTION_TIMEOUT=4"}
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

			// Monitor the container process for early exit
			cmdDone := make(chan error, 1)
			go func() {
				cmdDone <- cmd.Wait()
			}()

			// Poll until the container's RIE is accepting TCP connections.
			// We only do TCP dialing here — NOT HTTP requests — to avoid
			// sending multiple requests, which can result in failures.
			const pollInterval = 100 * time.Millisecond
			addr := fmt.Sprintf("127.0.0.1:%d", port)
			deadline := time.Now().Add(30 * time.Second)
			for time.Now().Before(deadline) {
				select {
				case err := <-cmdDone:
					<-logDone
					require.Failf(t, "container exited before becoming ready", "exit error: %v\nlogs:\n%s", err, logBuf.String())
				default:
				}
				conn, dialErr := net.DialTimeout("tcp", addr, pollInterval)
				if dialErr == nil {
					conn.Close()
					break
				}
				time.Sleep(pollInterval)
			}

			// Give the RIE a moment to fully initialize its HTTP handler after
			// the TCP listener is up.
			time.Sleep(500 * time.Millisecond)

			client := &http.Client{Timeout: 10 * time.Second}
			invokeURL := fmt.Sprintf("http://127.0.0.1:%d/2015-03-31/functions/function/invocations", port)
			resp, err := client.Post(invokeURL, "application/json", strings.NewReader("{}"))
			require.NoError(t, err)
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Task timed out after 4.00 seconds", string(body))

			_ = cmd.Process.Kill()
			<-cmdDone
			<-logDone

			logs := logBuf.String()
			t.Logf("stdout:\n%s", logs)
			opts.assertLogs(t, logs)
		})
	}
}
