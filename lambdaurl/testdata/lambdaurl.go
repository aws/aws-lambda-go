package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambdaurl"
)

const content = `<!DOCTYPE HTML>
<html>
<body>
Hello World!
</body>
</html>
`

func main() {
	lambdaurl.Start(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, strings.NewReader(content))
	}),
		lambdaurl.WithDetectContentType(true),
	)
}
