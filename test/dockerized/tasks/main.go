// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func dispatch(ctx context.Context, _ json.RawMessage) (interface{}, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no lambda context on invoke context")
	}

	switch handler := os.Getenv("_HANDLER"); handler {
	case "w3c.getW3c":
		return lc.W3C(), nil

	case "w3c.getW3cAndSource":
		_, hasW3c := lc.ClientContext.Custom["w3c"]
		return map[string]interface{}{
			"w3c":                 lc.W3C(),
			"clientContextCustom": lc.ClientContext.Custom,
			"clientContextHasW3c": hasW3c,
		}, nil

	case "w3c.echoClientContext":
		return map[string]interface{}{
			"custom": lc.ClientContext.Custom,
			"env":    lc.ClientContext.Env,
		}, nil

	default:
		return nil, fmt.Errorf("unknown handler: %q", handler)
	}
}

func main() {
	lambda.Start(dispatch)
}
