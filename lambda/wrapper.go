// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Specify the noexecwrapper build tag to remove the wrapper tampoline from
// this library if it is undesirable.
//go:build unix && !noexecwrapper
// +build unix,!noexecwrapper

package lambda

import (
	"log"
	"os"
	"syscall"
)

const awsLambdaExecWrapper = "AWS_LAMBDA_EXEC_WRAPPER"

func init() {
	// Honor the AWS_LAMBDA_EXEC_WRAPPER configuration at startup, trying to emulate
	// the behavior of managed runtimes, as this configuration is otherwise not applied
	// by provided runtimes (or go1.x).
	execAWSLambdaExecWrapper(os.Getenv, syscall.Exec)
}

// If AWS_LAMBDA_EXEC_WRAPPER is defined, replace the current process by spawning
// it with the current process' arguments (including the program name). If the call
// to syscall.Exec fails, this aborts the process with a fatal error.
func execAWSLambdaExecWrapper(
	getenv func(key string) string,
	sysExec func(argv0 string, argv []string, envv []string) error,
) {
	wrapper := getenv(awsLambdaExecWrapper)
	if wrapper == "" {
		return
	}

	// The AWS_LAMBDA_EXEC_WRAPPER variable is blanked before replacing the process
	// in order to avoid endlessly restarting the process.
	env := append(os.Environ(), awsLambdaExecWrapper+"=")
	if err := sysExec(wrapper, append([]string{wrapper}, os.Args...), env); err != nil {
		log.Fatalf("failed to sysExec() %s=%s: %v", awsLambdaExecWrapper, wrapper, err)
	}
}
