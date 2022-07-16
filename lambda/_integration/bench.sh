#!/bin/bash

set -euo pipefail

echo "\"$(od -N $((512 * 1024)) /dev/random | base64)\"" > data.json
echo "data payload for tests is: $(du -h data.json)"

trap "docker kill rie-bench" 1 2 3 6
bench () {
    echo "-------------------------------------------------"
    echo $@
    echo "-------------------------------------------------"
    docker run --name rie-bench --platform $5 --rm -d -p 9001:8080 -v "${PWD}/$1:/var/task/$1" --entrypoint aws-lambda-rie $4 $2 $3 
    sleep 2
    echo "ensuring healthy function before starting test"
    curl -s -XPOST http://localhost:9001/2015-03-31/functions/function/invocations -d '{"hello": "world"}' | jq
    echo "-------------------------------------------------"
    ab -p data.json -n 100 http://localhost:9001/2015-03-31/functions/function/invocations
    docker kill rie-bench
}

GOOS=linux GOARCH=amd64 go build                    -o handler       echo_handler.go
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o handler_norpc echo_handler.go
GOOS=linux GOARCH=arm64 go build -tags lambda.nprpc -o handler_arm   echo_handler.go
ls -lah handler*

bench handler_arm   /var/task/handler_arm  x              public.ecr.aws/lambda/provided linux/arm64/v8
bench handler       /var/task/handler      x              public.ecr.aws/lambda/provided linux/amd64
bench handler_norpc /var/runtime/bootstrap handler_norpc  public.ecr.aws/lambda/go       linux/amd64
bench handler       /var/runtime/bootstrap handler        public.ecr.aws/lambda/go       linux/amd64
