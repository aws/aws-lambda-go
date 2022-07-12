// Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.

//go:build lambda.norpc
// +build lambda.norpc

package lambda

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartNotInLambda(t *testing.T) {
	actual := "unexpected"
	logFatalf = func(format string, v ...interface{}) {
		actual = fmt.Sprintf(format, v...)
	}

	Start(func() error { return nil })
	assert.Equal(t, "expected AWS Lambda environment variables [AWS_LAMBDA_RUNTIME_API] are not defined", actual)
}
