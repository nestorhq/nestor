package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// CloudWatchAPI api
type CloudWatchAPI struct {
	resourceTags *ResourceTags
	client       *cloudwatch.CloudWatch
}

// CloudWatchInformation description of a CloudWatch
type CloudWatchInformation struct {
	functionName string
	arn          string
}

// NewCloudWatchAPI constructor
func NewCloudWatchAPI(session *session.Session, resourceTags *ResourceTags) (*CloudWatchAPI, error) {
	var api = CloudWatchAPI{resourceTags: resourceTags}
	// Create CloudWatch client
	api.client = cloudwatch.New(session)
	return &api, nil
}
