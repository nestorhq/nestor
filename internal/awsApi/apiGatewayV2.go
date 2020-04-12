package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
)

// APIGatewayV2API api
type APIGatewayV2API struct {
	resourceTags *ResourceTags
	client       *apigatewayv2.ApiGatewayV2
}

// APIGatewayV2Information description of a ApiGatewayV2
type APIGatewayV2Information struct {
	functionName string
	arn          string
}

// NewAPIGatewayV2API constructor
func NewAPIGatewayV2API(session *session.Session, resourceTags *ResourceTags) (*APIGatewayV2API, error) {
	var api = APIGatewayV2API{resourceTags: resourceTags}
	// Create ApiGatewayV2 client
	api.client = apigatewayv2.New(session)
	return &api, nil
}
