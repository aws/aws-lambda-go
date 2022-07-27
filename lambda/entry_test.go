// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ctxTestKey struct{}

func TestStartRuntimeAPIWithContext(t *testing.T) {
	server, _ := runtimeAPIServer("null", 1) // serve a single invoke, and then cause an internal error
	expected := "expected"
	actual := "unexpected"

	os.Setenv("AWS_LAMBDA_RUNTIME_API", strings.Split(server.URL, "://")[1])
	defer os.Unsetenv("AWS_LAMBDA_RUNTIME_API")
	logFatalf = func(format string, v ...interface{}) {}
	defer func() { logFatalf = log.Fatalf }()

	StartWithContext(context.WithValue(context.Background(), ctxTestKey{}, expected), func(ctx context.Context) error {
		actual, _ = ctx.Value(ctxTestKey{}).(string)
		return nil
	})

	assert.Equal(t, expected, actual)
}
