package main

import (
	"context"
	"github.com/segmentio/encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambdaurl"
)

func logLambdaRequest(ctx context.Context) {
	req, ok := lambdaurl.RequestFromContext(ctx)
	if ok {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		enc.Encode(req)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	logLambdaRequest(r.Context())
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, strings.NewReader(`<html><body>Hello World!</body></html>`))
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/", root)
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambdaurl.Start(http.DefaultServeMux)
	}
	http.ListenAndServe(":9001", nil)
}
