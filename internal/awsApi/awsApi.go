package awsApi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AwsApi struct {
	profileName string
	session     *session.Session
}

func NewAwsApi(profileName string) (*AwsApi, error) {
	var awsApi = AwsApi{profileName: profileName}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profileName,
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
	return &awsApi, nil
}
