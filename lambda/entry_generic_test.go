//go:build go1.18
// +build go1.18

// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved

package lambda

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartHandlerFunc(t *testing.T) {
	actual := "unexpected"
	logFatalf = func(format string, v ...interface{}) {
		actual = fmt.Sprintf(format, v...)
	}

	f := func(context.Context, any) (any, error) { return 1, nil }
	StartHandlerFunc(f)

	assert.Equal(t, "expected AWS Lambda environment variables [_LAMBDA_SERVER_PORT AWS_LAMBDA_RUNTIME_API] are not defined", actual)

	handlerType := reflect.TypeOf(f)

	handlerTakesContext, err := validateArguments(handlerType)
	assert.NoError(t, err)
	assert.True(t, handlerTakesContext)

	err = validateReturns(handlerType)
	assert.NoError(t, err)
}
