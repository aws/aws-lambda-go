// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package lambda

import (
	"log"
	"os"
)

const awsLambdaExecWrapper = "AWS_LAMBDA_EXEC_WRAPPER"

// execAWSLambdaExecWrapper applies the AWS_LAMBDA_EXEC_WRAPPER environment variable.
// If AWS_LAMBDA_EXEC_WRAPPER is defined, replace the current process by spawning
// it with the current process' arguments (including the program name). If the call
// to syscall.Exec fails, this aborts the process with a fatal error.
func execAWSLambdaExecWrapper(
	getenv func(key string) string,
	sysExec func(argv0 string, argv []string, envv []string) error,
	callbacks []func(),
) {
	wrapper := getenv(awsLambdaExecWrapper)
	if wrapper == "" {
		return
	}

	// Execute the provided callbacks before re-starting the process...
	for _, callback := range callbacks {
		callback()
	}

	// The AWS_LAMBDA_EXEC_WRAPPER variable is blanked before replacing the process
	// in order to avoid endlessly restarting the process.
	env := append(os.Environ(), awsLambdaExecWrapper+"=")
	if err := sysExec(wrapper, append([]string{wrapper}, os.Args...), env); err != nil {
		log.Fatalf("failed to sysExec() %s=%s: %v", awsLambdaExecWrapper, wrapper, err)
	}
}
