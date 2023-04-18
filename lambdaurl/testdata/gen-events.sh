#!/bin/bash
set -euo pipefail

url_id="${1}" # should be the lambda function url domain prefix for an echo function
region=${AWS_REGION:-us-west-2}
url="https://${url_id}.lambda-url.${region}.on.aws"
account_id=$(aws sts get-caller-identity --output text --query "Account")

redact () {
    #https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-unique-ids 
    sed "s/${url_id}/lambda-url-id/g" \
    | sed 's/A[A-Z][A-Z]A[A-Z1-9]*\([":]\)/iam-unique-id\1/g' \
    | sed "s/${account_id}/aws-account-id/g" \
    | jq '.headers |= (.["x-amz-security-token"] = "security-token" )' \
    | jq '.headers |= (.["x-forwarded-for"] = "127.0.0.1")' \
    | jq '.requestContext.authorizer |= (.["iam"] = {})' \
    | jq '.requestContext.http |= (.["sourceIp"] = "127.0.0.1")'
}

awscurl --service lambda --region $region \
    -X POST \
    -H 'Header1: h1' \
    -H 'Header2: h1,h2' \
    -H 'Header3: h3' \
    -H 'Cookie: foo=bar; hello=hello' \
    -H 'Content-Type: application/json' \
    -d '{"hello": "world"}' \
    "$url/hello?hello=world&foo=bar" \
    | redact \
    | tee function-url-request-with-headers-and-cookies-and-text-body.json \
    | jq

awscurl --service lambda --region $region \
    -X POST \
    -d '<idk/>' \
    "$url" \
    | redact \
    | tee function-url-domain-only-request-with-base64-encoded-body.json \
    | jq

awscurl --service lambda --region $region \
    -X GET \
    "$url" \
    | redact \
    | tee function-url-domain-only-get-request.json \
    | jq

awscurl --service lambda --region $region \
    -X GET \
    "$url/" \
    | redact \
    | tee function-url-domain-only-get-request-trailing-slash.json \
    | jq

