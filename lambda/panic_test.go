package lambda

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/stretchr/testify/assert"
)

func assertPanicMessage(t *testing.T, panicFunc func(), expectedMessage string) {
	defer func() {
		if err := recover(); err != nil {
			panicInfo := getPanicInfo(err)
			assert.NotNil(t, panicInfo)
			assert.NotNil(t, panicInfo.Message)
			assert.Equal(t, expectedMessage, panicInfo.Message)
			assert.NotNil(t, panicInfo.StackTrace)
		}
	}()

	panicFunc()
	t.Errorf("Should have exited due to panic")
}

func TestPanicFormattingStringValue(t *testing.T) {
	assertPanicMessage(t, func() { panic("Panic time!") }, "Panic time!")
}

func TestPanicFormattingIntValue(t *testing.T) {
	assertPanicMessage(t, func() { panic(1234) }, "1234")
}

func TestPanicFormattingCustomError(t *testing.T) {
	customError := &CustomError{}
	assertPanicMessage(t, func() { panic(customError) }, customError.Error())
}

func TestPanicFormattingInvokeResponse_Error(t *testing.T) {
	ive := &messages.InvokeResponse_Error{Message: "message", Type: "type"}
	assertPanicMessage(t, func() { panic(ive) }, ive.Error())
}

func TestFormatFrame(t *testing.T) {
	var tests = []struct {
		inputPath     string
		inputLine     int32
		inputLabel    string
		expectedPath  string
		expectedLine  int32
		expectedLabel string
	}{
		{
			inputPath:     "/Volumes/Unix/workspace/LambdaGoLang/src/GoAmzn-Github-Aws-AwsLambdaGo/src/github.com/aws/aws-lambda-go/lambda/panic_test.go",
			inputLine:     42,
			inputLabel:    "github.com/aws/aws-lambda-go/lambda.printStack",
			expectedPath:  "github.com/aws/aws-lambda-go/lambda/panic_test.go",
			expectedLine:  42,
			expectedLabel: "printStack",
		},
		{
			inputPath:     "/home/user/src/pkg/sub/file.go",
			inputLine:     42,
			inputLabel:    "pkg/sub.Type.Method",
			expectedPath:  "pkg/sub/file.go",
			expectedLine:  42,
			expectedLabel: "Type.Method",
		},
		{
			inputPath:     "/home/user/src/pkg/sub/sub2/file.go",
			inputLine:     42,
			inputLabel:    "pkg/sub/sub2.Type.Method",
			expectedPath:  "pkg/sub/sub2/file.go",
			expectedLine:  42,
			expectedLabel: "Type.Method",
		},
		{
			inputPath:     "/home/user/src/pkg/file.go",
			inputLine:     101,
			inputLabel:    "pkg.Type.Method",
			expectedPath:  "pkg/file.go",
			expectedLine:  101,
			expectedLabel: "Type.Method",
		},
	}

	for _, test := range tests {
		inputFrame := runtime.Frame{
			File:     test.inputPath,
			Line:     int(test.inputLine),
			Function: test.inputLabel,
		}

		actual := formatFrame(inputFrame)
		assert.Equal(t, test.expectedPath, actual.Path)
		assert.Equal(t, test.expectedLine, actual.Line)
		assert.Equal(t, test.expectedLabel, actual.Label)
	}
}

func TestRuntimeStackTrace(t *testing.T) {
	// implementing the test in the inner function to simulate an
	// additional stack frame that would exist in real life due to the
	// defer function.
	testRuntimeStackTrace(t)
}

func testRuntimeStackTrace(t *testing.T) {
	panicInfo := getPanicInfo("Panic time!")

	assert.NotNil(t, panicInfo)
	assert.NotNil(t, panicInfo.StackTrace)
	assert.True(t, len(panicInfo.StackTrace) > 0)

	packagePath, err := getPackagePath()
	assert.NoError(t, err)

	frame := panicInfo.StackTrace[0]

	assert.Equal(t, packagePath+"/panic_test.go", frame.Path)
	assert.True(t, frame.Line > 0)
	assert.Equal(t, "testRuntimeStackTrace", frame.Label)
}

func getPackagePath() (string, error) {
	fullPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var paths []string
	if runtime.GOOS == "windows" {
		paths = strings.Split(fullPath, "\\")
	} else {
		paths = strings.Split(fullPath, "/")
	}

	// The frame.Path will only contain the last 5 directories if there are more than 5 directories.
	if len(paths) >= 5 {
		paths = paths[len(paths)-4:]
	}
	return strings.Join(paths, "/"), nil
}
