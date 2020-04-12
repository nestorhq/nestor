package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaAPI api
type LambdaAPI struct {
	resourceTags *ResourceTags
	client       *lambda.Lambda
}

// LambdaInformation description of a lambda
type LambdaInformation struct {
	functionName string
	arn          string
}

// NewLambdaAPI constructor
func NewLambdaAPI(session *session.Session, resourceTags *ResourceTags) (*LambdaAPI, error) {
	var api = LambdaAPI{resourceTags: resourceTags}
	// Create Lambda client
	api.client = lambda.New(session)
	return &api, nil
}
