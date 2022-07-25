#!/bin/bash

set -euo pipefail

echo "\"$(od -N $((512 * 1024)) /dev/random | base64)\"" > data.json
echo "data payload for tests is: $(du -h data.json)"

trap "docker kill rie-bench" 1 2 3 6
bench () {
    local handler_exe=$1
    local entrypoint=$2
    local image=$3
    echo "-------------------------------------------------"
    echo $@
    echo "-------------------------------------------------"
    docker run --name rie-bench --rm -d -p 9001:8080 -v "${handler_exe}:/var/task/bootstrap" --entrypoint aws-lambda-rie ${image} ${entrypoint} bootstrap
    sleep 2
    echo "ensuring healthy function before starting test"
    curl -s -XPOST http://localhost:9001/2015-03-31/functions/function/invocations -d '{"hello": "world"}' | jq
    echo "-------------------------------------------------"
    ab -p data.json -n 100 http://localhost:9001/2015-03-31/functions/function/invocations
    docker kill rie-bench
}

GOOS=linux GOARCH=amd64 go build                    -o handler       echo_handler.go
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o handler_norpc echo_handler.go
ls -lah handler*

bench "$(pwd)/handler_norpc" /var/task/bootstrap    public.ecr.aws/lambda/provided:alami 
bench "$(pwd)/handler_norpc" /var/runtime/bootstrap public.ecr.aws/lambda/go             
bench "$(pwd)/handler"       /var/runtime/bootstrap public.ecr.aws/lambda/go             
