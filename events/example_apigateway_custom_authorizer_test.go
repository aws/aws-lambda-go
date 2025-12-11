package events_test

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// This is a simple TOKEN authorizer example to demonstrate how to use an authorization
// token to allow or deny a request. In this example, the caller named "user" is allowed to invoke
// a request if the client-supplied token value is "allow". The caller is not allowed to invoke
// the request if the token value is "deny". If the token value is "Unauthorized", the function
// returns the "Unauthorized" error with an HTTP status code of 401. For any other token value,
// the authorizer returns an "Invalid token" error.
func ExampleAPIGatewayCustomAuthorizerRequest() {
	lambda.Start(func(ctx context.Context, event *events.APIGatewayCustomAuthorizerRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
		token := event.AuthorizationToken
		switch strings.ToLower(token) {
		case "allow":
			return generatePolicy("user", "Allow", event.MethodArn), nil
		case "deny":
			return generatePolicy("user", "Deny", event.MethodArn), nil
		case "unauthorized":
			return nil, errors.New("Unauthorized")
		default:
			return nil, errors.New("Error: Invalid token")
		}
	})
}

func generatePolicy(principalID, effect, resource string) *events.APIGatewayCustomAuthorizerResponse {
	authResponse := &events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	authResponse.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	}
	return authResponse
}
