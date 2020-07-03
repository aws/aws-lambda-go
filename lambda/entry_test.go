// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/lambda/messages"
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

func TestStartRPCWithContext(t *testing.T) {
	expected := "expected"
	actual := "unexpected"
	port := getFreeTCPPort()
	os.Setenv("_LAMBDA_SERVER_PORT", fmt.Sprintf("%d", port))
	defer os.Unsetenv("_LAMBDA_SERVER_PORT")
	go StartWithContext(context.WithValue(context.Background(), ctxTestKey{}, expected), func(ctx context.Context) error {
		actual, _ = ctx.Value(ctxTestKey{}).(string)
		return nil
	})

	var client *rpc.Client
	var pingResponse messages.PingResponse
	var invokeResponse messages.InvokeResponse
	var err error
	for {
		client, err = rpc.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			continue
		}
		break
	}
	for {
		if err := client.Call("Function.Ping", &messages.PingRequest{}, &pingResponse); err != nil {
			continue
		}
		break
	}
	if err := client.Call("Function.Invoke", &messages.InvokeRequest{}, &invokeResponse); err != nil {
		t.Logf("error invoking function: %v", err)
	}

	assert.Equal(t, expected, actual)
}

func getFreeTCPPort() int {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal("getFreeTCPPort failed: ", err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

func TestStartNotInLambda(t *testing.T) {
	actual := "unexpected"
	logFatalf = func(format string, v ...interface{}) {
		actual = fmt.Sprintf(format, v...)
	}

	Start(func() error { return nil })
	assert.Equal(t, "expected AWS Lambda environment variables [_LAMBDA_SERVER_PORT AWS_LAMBDA_RUNTIME_API] are not defined", actual)
}
