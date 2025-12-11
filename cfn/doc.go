// Package cfn provides helpers for implementing AWS CloudFormation custom resources.
//
// CloudFormation custom resources allow you to write custom provisioning logic that CloudFormation
// runs when you create, update, or delete stacks. This package handles the response protocol,
// making it easier to implement custom resource handlers.
//
// The LambdaWrap helper catches errors and ensures proper responses are sent to CloudFormation's
// pre-signed URL, preventing stack operations from hanging.
//
// See https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html
package cfn
