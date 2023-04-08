package handlertrace

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	ctx := context.Background()

	requestCall := 0
	responseCall := 0

	existedContext := context.WithValue(ctx, handlerTraceKey{}, HandlerTrace{
		RequestEvent: func(ctx context.Context, event interface{}) {
			requestCall += 1
		},
		ResponseEvent: func(ctx context.Context, event interface{}) {
			responseCall += 1
		},
	})

	trace := HandlerTrace{
		RequestEvent: func(ctx context.Context, event interface{}) {
			requestCall += 1
		},
		ResponseEvent: func(ctx context.Context, event interface{}) {
			responseCall += 1
		},
	}

	ctxWithTrace := NewContext(existedContext, trace)
	traceFromCtx := FromContext(ctxWithTrace)

	traceFromCtx.RequestEvent(ctxWithTrace, nil)
	assert.Equal(t, requestCall, 2)

	traceFromCtx.ResponseEvent(ctxWithTrace, nil)
	fmt.Println(responseCall)
	assert.Equal(t, responseCall, 2)
}
