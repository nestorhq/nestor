package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoAPI api
type CognitoAPI struct {
	resourceTags *ResourceTags
	client       *cognitoidentityprovider.CognitoIdentityProvider
}

// UserPoolInformation description of a user pool
type UserPoolInformation struct {
	ID  string
	arn string
}

// NewCognitoAPI constructor
func NewCognitoAPI(session *session.Session, resourceTags *ResourceTags, cognitoRegion string) (*CognitoAPI, error) {
	var api = CognitoAPI{resourceTags: resourceTags}
	// Create CognitoIdentityProvider client
	api.client = cognitoidentityprovider.New(session, aws.NewConfig().WithRegion(cognitoRegion))
	return &api, nil
}

func (api *CognitoAPI) getUserPoolInformationAndTags(userPoolID string) (*UserPoolInformation, map[string]*string, error) {
	var result = UserPoolInformation{ID: userPoolID}
	info, err := api.client.DescribeUserPool(&cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: &userPoolID,
	})
	if err != nil {
		return &result, nil, err
	}
	result.arn = *info.UserPool.Arn
	return &result, info.UserPool.UserPoolTags, nil
}

func (api *CognitoAPI) findUserPoolByName(userPoolName string) (string, error) {
	var nextToken = ""
	for {
		var input = cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int64(32),
		}
		if len(nextToken) > 0 {
			input.NextToken = &nextToken
		}

		listPools, err := api.client.ListUserPools(&input)
		if err != nil {
			return "", err
		}
		for _, userPool := range listPools.UserPools {
			if userPool.Name == &userPoolName {
				return *userPool.Id, nil
			}
		}
		if listPools.NextToken != nil {
			nextToken = *listPools.NextToken
		} else {
			// the pool was not found
			return "", nil
		}
	}
}

// doc at:
// https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentityprovider/#CreateUserPoolInput
func (api *CognitoAPI) doCreateUserPool(userPoolName string) (*UserPoolInformation, error) {
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
		UsernameAttributes: []*string{
			aws.String("email"),
		},
		UsernameConfiguration: &cognitoidentityprovider.UsernameConfigurationType{
			CaseSensitive: aws.Bool(false),
		},
		VerificationMessageTemplate: &cognitoidentityprovider.VerificationMessageTemplateType{
			DefaultEmailOption: aws.String("CONFIRM_WITH_CODE"),
			EmailMessage:       aws.String("Your verification code is {####}"),
			EmailSubject:       aws.String("Your verification code is {####}"),
			SmsMessage:         aws.String("Your verification code is {####}"),
		},
	}

	result, err := api.client.CreateUserPool(input)
	if err != nil {
		fmt.Println("Got error calling CreateUserPool:")
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Printf("user pool: %v\n", result)
	fmt.Println("Created the user pool", userPoolName)
	return &UserPoolInformation{
		ID:  *result.UserPool.Id,
		arn: *result.UserPool.Arn,
	}, nil
}

func (api *CognitoAPI) createUserPool(userPoolName string) (*UserPoolInformation, error) {
	id, err := api.findUserPoolByName(userPoolName)
	if err != nil {
		fmt.Println("Got error calling findUserPoolByName:")
		fmt.Println(err.Error())
		return nil, err
	}

	if id != "" {
		info, tags, err := api.getUserPoolInformationAndTags(id)
		if err != nil {
			return nil, err
		}
		// check tags
		err2 := api.resourceTags.checkTags(tags)
		if err2 != nil {
			return nil, err2
		}
		return info, nil
	}
	return api.doCreateUserPool(userPoolName)
}
