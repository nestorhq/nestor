package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/nestorhq/nestor/internal/reporter"
)

// AwsAPI api to work on AWS
type AwsAPI struct {
	profileName  string
	session      *session.Session
	resourceTags *ResourceTags
	dynamoDbAPI  *DynamoDbAPI
	cognitoAPI   *CognitoAPI
}

// NewAwsAPI constructor
func NewAwsAPI(profileName string, resourceTags *ResourceTags, region string, cognitoRegion string) (*AwsAPI, error) {
	var r0 = reporter.NewReporterM(
		reporter.NewMessage("Aws API initialization").
			WithArg("appName", resourceTags.appName).
			WithArg("environment", resourceTags.environment).
			WithArg("nestorVersion", resourceTags.nestorVersion).
			WithArg("profileName", profileName).
			WithArg("region", region).
			WithArg("cognitoRegion", cognitoRegion))

	t0 := r0.Start()

	var awsAPI = AwsAPI{profileName: profileName, resourceTags: resourceTags}

	t1 := t0.Sub("create AWS session")
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		t1.Fail(err)
		return nil, err
	}
	awsAPI.session = sess
	t1.Okr(map[string]string{
		"region": *sess.Config.Region,
	})

	//	fmt.Printf("region: %v\n", sess.Config.Endpoint)

	// Create a STS client from just a session.
	t2 := t0.Sub("create sts client and sts.GetCallerIdentityInput")
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		t2.Fail(err)
		return nil, err
	}
	// fmt.Println(result)
	t2.Okr(map[string]string{
		"account": *result.Account,
		"arn":     *result.Arn,
		"userId":  *result.UserId,
	})

	t3 := t0.Sub("create dynamoDb API")
	// initialize different AWS Apis
	awsAPI.dynamoDbAPI, err = NewDynamoDbAPI(sess, resourceTags)
	if err != nil {
		t3.Fail(err)
		return nil, err
	}
	t4 := t0.Sub("create Cognito API")
	awsAPI.cognitoAPI, err = NewCognitoAPI(sess, resourceTags, cognitoRegion)
	if err != nil {
		t4.Fail(err)
		return nil, err
	}
	r0.Ok()
	return &awsAPI, nil
}

// CreateMonoTable create a mongoDb table following the mono-table schema
func (api *AwsAPI) CreateMonoTable(tableName string) {
	api.dynamoDbAPI.createMonoTable(tableName)
}

// CreateUserPool create a user pool
func (api *AwsAPI) CreateUserPool(userPoolName string) (*UserPoolInformation, error) {
	return api.cognitoAPI.createUserPool(userPoolName)
}
