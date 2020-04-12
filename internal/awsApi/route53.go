package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Route53API api
type Route53API struct {
	resourceTags *ResourceTags
	client       *route53.Route53
}

// Route53Information description of a Route53
type Route53Information struct {
	functionName string
	arn          string
}

// NewRoute53API constructor
func NewRoute53API(session *session.Session, resourceTags *ResourceTags) (*Route53API, error) {
	var api = Route53API{resourceTags: resourceTags}
	// Create Route53 client
	api.client = route53.New(session)
	return &api, nil
}
