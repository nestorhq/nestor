package awsapi

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoAPI api
type CognitoAPI struct {
	client *cognitoidentityprovider.CognitoIdentityProvider
}

// NewCognitoAPI constructor
func NewCognitoAPI(session *session.Session, cognitoRegion string) (*CognitoAPI, error) {
	var api = CognitoAPI{}
	// Create CognitoIdentityProvider client
	api.client = cognitoidentityprovider.New(session, aws.NewConfig().WithRegion(cognitoRegion))
	return &api, nil
}

// doc at:
// https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentityprovider/#CreateUserPoolInput
func (api *CognitoAPI) createUserPool(userPoolName string) {
	input := &cognitoidentityprovider.CreateUserPoolInput{
		PoolName: &userPoolName,
		AdminCreateUserConfig: &cognitoidentityprovider.AdminCreateUserConfigType{
			AllowAdminCreateUserOnly: aws.Bool(false),
			InviteMessageTemplate: &cognitoidentityprovider.MessageTemplateType{
				SMSMessage:   aws.String("SMSMessage {####} {username}"),
				EmailMessage: aws.String("EmailMessage {####} {username}"),
				EmailSubject: aws.String("EmailSubject {####} {username}"),
			},
		},
		DeviceConfiguration: &cognitoidentityprovider.DeviceConfigurationType{
			ChallengeRequiredOnNewDevice:     aws.Bool(true),
			DeviceOnlyRememberedOnUserPrompt: aws.Bool(true),
		},
		EmailConfiguration: &cognitoidentityprovider.EmailConfigurationType{
			// TODO: use own SES account
			EmailSendingAccount: aws.String("COGNITO_DEFAULT"),
		},
	}

	result, err := api.client.CreateUserPool(input)
	if err != nil {
		fmt.Println("Got error calling CreateUserPool:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("user pool: %v\n", result)
	fmt.Println("Created the user pool", userPoolName)
}
