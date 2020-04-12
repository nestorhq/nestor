package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

// EventBridgeAPI api
type EventBridgeAPI struct {
	resourceTags *ResourceTags
	client       *eventbridge.EventBridge
}

// EventBridgeInformation description of a EventBridge
type EventBridgeInformation struct {
	functionName string
	arn          string
}

// NewEventBridgeAPI constructor
func NewEventBridgeAPI(session *session.Session, resourceTags *ResourceTags) (*EventBridgeAPI, error) {
	var api = EventBridgeAPI{resourceTags: resourceTags}
	// Create EventBridge client
	api.client = eventbridge.New(session)
	return &api, nil
}
