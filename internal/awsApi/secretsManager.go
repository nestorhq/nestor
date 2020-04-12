package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretsManagerAPI api
type SecretsManagerAPI struct {
	resourceTags *ResourceTags
	client       *secretsmanager.SecretsManager
}

// SecretsManagerInformation description of a SecretsManager
type SecretsManagerInformation struct {
	functionName string
	arn          string
}

// NewSecretsManagerAPI constructor
func NewSecretsManagerAPI(session *session.Session, resourceTags *ResourceTags) (*SecretsManagerAPI, error) {
	var api = SecretsManagerAPI{resourceTags: resourceTags}
	// Create SecretsManager client
	api.client = secretsmanager.New(session)
	return &api, nil
}
