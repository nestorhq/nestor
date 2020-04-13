package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/nestorhq/nestor/internal/reporter"
)

// AwsAPI api to work on AWS
type AwsAPI struct {
	profileName    string
	session        *session.Session
	resourceTags   *ResourceTags
	dynamoDbAPI    *DynamoDbAPI
	cognitoAPI     *CognitoAPI
	lambdaAPI      *LambdaAPI
	eventBridgeAPI *EventBridgeAPI
	s3API          *S3API
}

// NewAwsAPI constructor
func NewAwsAPI(profileName string, resourceTags *ResourceTags, region string, cognitoRegion string, t *reporter.Task) (*AwsAPI, error) {
	t0 := t.Sub("create Aws API")
	t0.LogM(
		reporter.NewMessage("Aws API initialization").
			WithArg("appName", resourceTags.appName).
			WithArg("environment", resourceTags.environment).
			WithArg("nestorVersion", resourceTags.nestorVersion).
			WithArg("profileName", profileName).
			WithArg("region", region).
			WithArg("cognitoRegion", cognitoRegion))

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
	t5 := t0.Sub("create Lambda API")
	awsAPI.lambdaAPI, err = NewLambdaAPI(sess, resourceTags)
	if err != nil {
		t5.Fail(err)
		return nil, err
	}

	t6 := t0.Sub("create EventBridge API")
	awsAPI.eventBridgeAPI, err = NewEventBridgeAPI(sess, resourceTags)
	if err != nil {
		t6.Fail(err)
		return nil, err
	}

	t7 := t0.Sub("create S3 API")
	awsAPI.s3API, err = NewS3API(sess, resourceTags)
	if err != nil {
		t7.Fail(err)
		return nil, err
	}

	t0.Ok()
	return &awsAPI, nil
}

// CreateUserPool create a user pool
func (api *AwsAPI) CreateUserPool(userPoolName string, nestorID string, t *reporter.Task) (*UserPoolInformation, error) {
	t0 := t.SubM(
		reporter.NewMessage("Aws API: CreateUserPool").
			WithArg("userPoolName", userPoolName))

	up, err := api.cognitoAPI.createUserPool(userPoolName, nestorID, t0)
	if err != nil {
		t0.Fail(err)
	} else {
		t0.Okr(map[string]string{
			"id":  up.ID,
			"arn": up.arn,
		})
	}

	return up, err
}

// CreateMonoTable create a mongoDb table following the mono-table schema
func (api *AwsAPI) CreateMonoTable(tableName string, nestorID string, t *reporter.Task) (*TableInformation, error) {
	t0 := t.SubM(
		reporter.NewMessage("Aws API: CreateMonoTable").
			WithArg("tableName", tableName))

	res, error := api.dynamoDbAPI.createMonoTable(tableName, nestorID, t0)
	if error != nil {
		t0.Fail(error)
	}
	return res, error
}

// CreateEventBus create event bus
func (api *AwsAPI) CreateEventBus(eventBusName string, nestorID string, t *reporter.Task) (*EventBusInformation, error) {
	t0 := t.SubM(
		reporter.NewMessage("Aws API: CreateEventBus").
			WithArg("eventBusName", eventBusName))

	res, error := api.eventBridgeAPI.createEventBus(eventBusName, nestorID, t0)
	if error != nil {
		t0.Fail(error)
	}
	return res, error
}

// CreateBucket create bucket
func (api *AwsAPI) CreateBucket(bucketName string, nestorID string, t *reporter.Task) (*S3Information, error) {
	t0 := t.SubM(
		reporter.NewMessage("Aws API: CreateBucket").
			WithArg("bucketName", bucketName))

	res, error := api.s3API.createBucket(bucketName, nestorID, t0)
	if error != nil {
		t0.Fail(error)
	}
	return res, error
}
