package awsapi

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoAPI api
type CognitoAPI struct {
	client *cognitoidentityprovider.CognitoIdentityProvider
}

// NewCognitoAPI constructor
func NewCognitoAPI(session *session.Session) (*CognitoAPI, error) {
	var api = CognitoAPI{}
	// Create CognitoIdentityProvider client
	api.client = cognitoidentityprovider.New(session)
	return &api, nil
}
