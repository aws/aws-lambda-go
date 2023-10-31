// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package lambda

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecAwsLambdaExecWrapperNotSet(t *testing.T) {
	var called bool
	callback := func() { called = true }

	exec, execCalled := mockExec(t, "<nope>")
	execAWSLambdaExecWrapper(
		mockedGetenv(t, ""),
		exec,
		[]func(){callback},
	)
	require.False(t, *execCalled)
	require.False(t, called)
}

func TestExecAwsLambdaExecWrapperSet(t *testing.T) {
	var called bool
	callback := func() { called = true }

	wrapper := "/path/to/wrapper/entry/point"
	exec, execCalled := mockExec(t, wrapper)
	execAWSLambdaExecWrapper(
		mockedGetenv(t, wrapper),
		exec,
		[]func(){callback},
	)
	require.True(t, *execCalled)
	require.True(t, called)
}

func mockExec(t *testing.T, value string) (mock func(string, []string, []string) error, called *bool) {
	mock = func(argv0 string, argv []string, envv []string) error {
		*called = true
		require.Equal(t, value, argv0)
		require.Equal(t, append([]string{value}, os.Args...), argv)
		require.Equal(t, awsLambdaExecWrapper+"=", envv[len(envv)-1])
		return nil
	}
	called = ptrTo(false)
	return
}

func mockedGetenv(t *testing.T, value string) func(string) string {
	return func(key string) string {
		require.Equal(t, awsLambdaExecWrapper, key)
		return value
	}
}

func ptrTo(val bool) *bool {
	return &val
}
