package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/nestorhq/nestor/internal/reporter"
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
		// look for the user pool given by name
		for _, userPool := range listPools.UserPools {
			if *userPool.Name == userPoolName {
				return *userPool.Id, nil
			}
		}
		// check if we have to paginate
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
func (api *CognitoAPI) doCreateUserPool(userPoolName string, nestorID string, task *reporter.Task) (*UserPoolInformation, error) {
	t0 := task.SubM(reporter.NewMessage("cognitoidentityprovider.CreateUserPoolInput").WithArg("userPoolName", userPoolName))
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
			{
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
			{
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
			{
				AttributeDataType:      aws.String("Boolean"),
				DeveloperOnlyAttribute: aws.Bool(false),
				Required:               aws.Bool(false),
				Name:                   aws.String("email_verified"),
				Mutable:                aws.Bool(true),
			},
			{
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
		UserPoolTags: aws.StringMap(api.resourceTags.getTagsAsMapWithID(nestorID)),
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
		t0.Fail(err)
		return nil, err
	}
	t0.Okr(map[string]string{
		"ID":  *result.UserPool.Id,
		"arn": *result.UserPool.Arn,
	})

	return &UserPoolInformation{
		ID:  *result.UserPool.Id,
		arn: *result.UserPool.Arn,
	}, nil
}

func (api *CognitoAPI) createUserPool(userPoolName string, nestorID string, t *reporter.Task) (*UserPoolInformation, error) {
	t0 := t.SubM(reporter.NewMessage("findUserPoolByName").WithArg("userPoolName", userPoolName))
	id, err := api.findUserPoolByName(userPoolName)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}

	if id != "" {
		t0.Log("user pool already exists")
		t0.Ok()
		t1 := t.SubM(reporter.NewMessage("getUserPoolInformationAndTags").WithArg("id", id))
		info, tags, err := api.getUserPoolInformationAndTags(id)
		if err != nil {
			t1.Fail(err)
			return nil, err
		}
		// check tags
		t2 := t.SubM(reporter.NewMessage("checkTags").WithArgs(tags))
		err2 := api.resourceTags.checkTags(tags, nestorID)
		if err2 != nil {
			t2.Fail(err2)
			return nil, err2
		}
		t2.Ok()
		return info, nil
	}
	t0.Log("user pool does not exist")
	t0.Ok()
	t3 := t.SubM(reporter.NewMessage("doCreateUserPool").WithArg("userPoolName", userPoolName))
	return api.doCreateUserPool(userPoolName, nestorID, t3)
}
