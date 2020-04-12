package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3API api
type S3API struct {
	resourceTags *ResourceTags
	client       *s3.S3
}

// S3Information description of a S3
type S3Information struct {
	functionName string
	arn          string
}

// NewS3API constructor
func NewS3API(session *session.Session, resourceTags *ResourceTags) (*S3API, error) {
	var api = S3API{resourceTags: resourceTags}
	// Create S3 client
	api.client = s3.New(session)
	return &api, nil
}
