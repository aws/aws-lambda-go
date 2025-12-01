//go:build go1.22
// +build go1.22

// Copyright 2025 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeAPILoopWithConcurrency(t *testing.T) {
	nInvokes := 100
	concurrency := 5

	metadata := make([]eventMetadata, nInvokes)
	for i := range nInvokes {
		m := defaultInvokeMetadata()
		m.requestID = fmt.Sprintf("request-%d", i)
		metadata[i] = m
	}

	ts, record := runtimeAPIServer(``, nInvokes, metadata...)
	defer ts.Close()

	active := atomic.Int32{}
	maxActive := atomic.Int32{}
	handler := NewHandler(func(ctx context.Context) (string, error) {
		activeNow := active.Add(1)
		defer active.Add(-1)
		for pr := maxActive.Load(); activeNow > pr; pr = maxActive.Load() {
			if maxActive.CompareAndSwap(pr, activeNow) {
				break
			}
		}
		lc, _ := lambdacontext.FromContext(ctx)
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		switch lc.AwsRequestID[len(lc.AwsRequestID)-1:] {
		case "6", "7":
			return "", fmt.Errorf("error-%s", lc.AwsRequestID)
		default:
			return lc.AwsRequestID, nil
		}
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoopWithConcurrency(endpoint, handler, concurrency), expectedError)
	assert.GreaterOrEqual(t, record.nGets, nInvokes+1)
	assert.Equal(t, nInvokes, record.nPosts)
	assert.Equal(t, int32(concurrency), maxActive.Load())
	responses := make(map[string]int)
	for _, response := range record.responses {
		responses[string(response)]++
	}
	assert.Len(t, responses, nInvokes)
	for response, count := range responses {
		assert.Equal(t, 1, count, "response %s seen %d times", response, count)
	}
	for i := range nInvokes {
		switch i % 10 {
		case 6, 7:
			assert.Contains(t, responses, fmt.Sprintf(`{"errorMessage":"error-request-%d","errorType":"errorString"}`, i))
		default:
			assert.Contains(t, responses, fmt.Sprintf(`"request-%d"`, i))
		}
	}
}

func TestRuntimeAPILoopSingleConcurrency(t *testing.T) {
	nInvokes := 10

	ts, record := runtimeAPIServer(``, nInvokes)
	defer ts.Close()

	var counter atomic.Int32
	handler := NewHandler(func(ctx context.Context) (string, error) {
		counter.Add(1)
		return "Hello!", nil
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	expectedError := fmt.Sprintf("failed to GET http://%s/2018-06-01/runtime/invocation/next: got unexpected status code: 410", endpoint)
	assert.EqualError(t, startRuntimeAPILoopWithConcurrency(endpoint, handler, 1), expectedError)
	assert.Equal(t, nInvokes+1, record.nGets)
	assert.Equal(t, nInvokes, record.nPosts)
	assert.Equal(t, int32(nInvokes), counter.Load())
}

func TestRuntimeAPILoopWithConcurrencyPanic(t *testing.T) {
	concurrency := 3

	ts, record := runtimeAPIServer(``, 100)
	defer ts.Close()

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	var counter atomic.Int32
	handler := NewHandler(func() error {
		n := counter.Add(1)
		time.Sleep(time.Duration(n) * 10 * time.Millisecond)
		panic(fmt.Errorf("panic %d", n))
	})
	endpoint := strings.Split(ts.URL, "://")[1]
	err := startRuntimeAPILoopWithConcurrency(endpoint, handler, concurrency)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "calling the handler function resulted in a panic, the process should exit")
	assert.Equal(t, concurrency, record.nGets)
	assert.Equal(t, concurrency, record.nPosts)
	assert.Equal(t, int32(concurrency), counter.Load())
	assert.Contains(t, string(record.responses[0]), "panic 1")
	logs := logBuf.String()
	idx1 := strings.Index(logs, "panic 1")
	idx2 := strings.Index(logs, "panic 2")
	idx3 := strings.Index(logs, "panic 3")
	assert.Greater(t, idx1, -1)
	assert.Greater(t, idx2, idx1)
	assert.Greater(t, idx3, idx2)
}
