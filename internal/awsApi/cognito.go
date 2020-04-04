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
	resourceTags *ResourceTags
	client       *cognitoidentityprovider.CognitoIdentityProvider
}

// NewCognitoAPI constructor
func NewCognitoAPI(session *session.Session, resourceTags *ResourceTags, cognitoRegion string) (*CognitoAPI, error) {
	var api = CognitoAPI{resourceTags: resourceTags}
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
		EmailVerificationMessage: aws.String("Your verification code is {####}"),
		EmailVerificationSubject: aws.String("Your verification code is {####}"),
		LambdaConfig:             &cognitoidentityprovider.LambdaConfigType{},
		MfaConfiguration:         aws.String("OFF"),
		Policies: &cognitoidentityprovider.UserPoolPolicyType{
			PasswordPolicy: &cognitoidentityprovider.PasswordPolicyType{
				MinimumLength:                 aws.Int64(8),
				RequireUppercase:              aws.Bool(false),
				RequireLowercase:              aws.Bool(false),
				RequireNumbers:                aws.Bool(true),
				RequireSymbols:                aws.Bool(true),
				TemporaryPasswordValidityDays: aws.Int64(7),
			},
		},
		Schema: []*cognitoidentityprovider.SchemaAttributeType{
			&cognitoidentityprovider.SchemaAttributeType{
				Name: aws.String("sub"),
				StringAttributeConstraints: &cognitoidentityprovider.StringAttributeConstraintsType{
					MinLength: aws.String("1"),
					MaxLength: aws.String("2048"),
				},
				DeveloperOnlyAttribute: aws.Bool(false),
				Required:               aws.Bool(true),
				AttributeDataType:      aws.String("String"),
				Mutable:                aws.Bool(false),
			},
			&cognitoidentityprovider.SchemaAttributeType{
				Name: aws.String("email"),
				StringAttributeConstraints: &cognitoidentityprovider.StringAttributeConstraintsType{
					MinLength: aws.String("0"),
					MaxLength: aws.String("2048"),
				},
				DeveloperOnlyAttribute: aws.Bool(false),
				Required:               aws.Bool(true),
				AttributeDataType:      aws.String("String"),
				Mutable:                aws.Bool(true),
			},
			&cognitoidentityprovider.SchemaAttributeType{
				AttributeDataType:      aws.String("Boolean"),
				DeveloperOnlyAttribute: aws.Bool(false),
				Required:               aws.Bool(false),
				Name:                   aws.String("email_verified"),
				Mutable:                aws.Bool(true),
			},
			&cognitoidentityprovider.SchemaAttributeType{
				Name: aws.String("updated_at"),
				NumberAttributeConstraints: &cognitoidentityprovider.NumberAttributeConstraintsType{
					MinValue: aws.String("0"),
				},
				DeveloperOnlyAttribute: aws.Bool(false),
				Required:               aws.Bool(false),
				AttributeDataType:      aws.String("Number"),
				Mutable:                aws.Bool(true),
			},
		},
		SmsAuthenticationMessage: aws.String("Your verification code is {####}"),
		SmsVerificationMessage:   aws.String("Your verification code is {####}"),
		UserPoolAddOns: &cognitoidentityprovider.UserPoolAddOnsType{
			AdvancedSecurityMode: aws.String("AUDIT"),
		},
		UserPoolTags: aws.StringMap(api.resourceTags.getTagsAsMap()),
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

/*
		Schema: &[]*cognitoidentityprovider.SchemaAttributeType{&{
      Name: "sub",
      StringAttributeConstraints: {
        MinLength: "1",
        MaxLength: "2048"
      },
      DeveloperOnlyAttribute: false,
      Required: true,
      AttributeDataType: "String",
      Mutable: false
    },
    &{
      Name: "email",
      StringAttributeConstraints: {
        MinLength: "0",
        MaxLength: "2048"
      },
      DeveloperOnlyAttribute: false,
      Required: true,
      AttributeDataType: "String",
      Mutable: true
    },
    &{
      AttributeDataType: "Boolean",
      DeveloperOnlyAttribute: false,
      Required: false,
      Name: "email_verified",
      Mutable: true
    },
    &{
      Name: "updated_at",
      NumberAttributeConstraints: {
        MinValue: "0"
      },
      DeveloperOnlyAttribute: false,
      Required: false,
      AttributeDataType: "Number",
      Mutable: true
    }},
	}

*/
