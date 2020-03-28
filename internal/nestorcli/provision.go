package nestorcli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/nestorhq/nestor/internal/config"
)

func CliProvision(environment string, nestorConfig *config.Config) {
	fmt.Printf("CliProvision: Environment is %s\n", environment)
	fmt.Printf("config: %v\n", nestorConfig)

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "sls",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("region: %v\n", sess.Config.Endpoint)

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
		return
	}

	fmt.Println(result)

}
