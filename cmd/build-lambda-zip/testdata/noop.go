// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

type BinaryHandler func(context.Context, []byte) ([]byte, error)
func (bh BinaryHandler) Invoke(ctx context.Context, req []byte) ([]byte, error) {
	return bh(ctx, req)
}

func noop (ctx context.Context, req []byte) ([]byte, error) {
	return req, nil
}

func main() {
	lambda.StartHandler(BinaryHandler(noop))
}
