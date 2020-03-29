package awsApi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AwsApi struct {
	profileName string
	session     *session.Session
	dynamoDbApi *DynamoDbApi
}

func NewAwsApi(profileName string, region string) (*AwsApi, error) {
	var awsApi = AwsApi{profileName: profileName}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		return nil, err
	}
	awsApi.session = sess

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
	awsApi.dynamoDbApi, err = NewDynameDbApi(sess)
	if err != nil {
		return nil, err
	}
	return &awsApi, nil
}

func (api *AwsApi) CreateMonoTable(tableName string) {
	api.dynamoDbApi.createMonoTable(tableName)
}
