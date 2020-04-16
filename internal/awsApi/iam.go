package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAMAPI api
type IAMAPI struct {
	resourceTags *ResourceTags
	client       *iam.IAM
}

// IAMInformation description of a IAM
type IAMInformation struct {
	functionName string
	arn          string
}

// RoleInformation  description of an IAM role
type RoleInformation struct {
	RoleArn string
}

// NewIAMAPI constructor
func NewIAMAPI(session *session.Session, resourceTags *ResourceTags) (*IAMAPI, error) {
	var api = IAMAPI{resourceTags: resourceTags}
	// Create IAM client
	api.client = iam.New(session)
	return &api, nil
}
