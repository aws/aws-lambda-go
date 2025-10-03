//go:build !go1.22
// +build !go1.22

// Copyright 2025 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
)

func startRuntimeAPILoop(api string, handler Handler) error {
	return doRuntimeAPILoop(context.Background(), newRuntimeAPIClient(api), newHandler(handler))
}
