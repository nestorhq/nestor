package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// AwsAPI api to work on AWS
type AwsAPI struct {
	profileName string
	session     *session.Session
	dynamoDbAPI *DynamoDbAPI
}

// NewAwsAPI constructor
func NewAwsAPI(profileName string, region string) (*AwsAPI, error) {
	var awsAPI = AwsAPI{profileName: profileName}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		return nil, err
	}
	awsAPI.session = sess

	//	fmt.Printf("region: %v\n", sess.Config.Endpoint)

	// Create a STS client from just a session.
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}
	fmt.Println(result)

	// initialize different AWS Apis
	awsAPI.dynamoDbAPI, err = NewDynamoDbAPI(sess)
	if err != nil {
		return nil, err
	}
	return &awsAPI, nil
}

// CreateMonoTable create a mongoDb table following the mono-table schema
func (api *AwsAPI) CreateMonoTable(tableName string) {
	api.dynamoDbAPI.createMonoTable(tableName)
}
