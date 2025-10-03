//go:build go1.22
// +build go1.22

// Copyright 2025 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func startRuntimeAPILoop(api string, handler Handler) error {
	h := newHandler(handler)
	client := newRuntimeAPIClient(api)
	concurrency := lambdacontext.MaxConcurrency()
	if concurrency <= 1 {
		return doRuntimeAPILoop(context.Background(), client, h)
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("no handlers run"))

	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for range concurrency {
		go func() {
			cancel(doRuntimeAPILoop(ctx, client, h))
			wg.Done()
		}()
	}
	wg.Wait()

	return context.Cause(ctx)
}
