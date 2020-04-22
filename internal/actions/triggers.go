package actions

import (
	"errors"

	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/reporter"
	"github.com/nestorhq/nestor/internal/resources"
)

// CreateTriggers processing for provision CLI command
// reference:
// https://docs.aws.amazon.com/AmazonS3/latest/dev/NotificationHowTo.html
func (actions *Actions) CreateTriggers(task *reporter.Task) error {
	// var appName = actions.nestorConfig.Application.Name
	// var environment = actions.environment
	var api = actions.api
	var nestorConfig = actions.nestorConfig
	var nestorResources = actions.nestorResources

	var err error
	var triggers = nestorConfig.Triggers
	// s3 triggers
	for _, s3CopyTrigger := range triggers.S3copy {
		s3Resource := nestorResources.FindResourceByID(s3CopyTrigger.BucketID)
		if s3Resource == nil {
			return errors.New("s3uploadTrigger: bucket not registered:" + s3CopyTrigger.BucketID)
		}
		bucketArn := s3Resource.GetAttribute(resources.AttArn)
		bucketName := s3Resource.GetAttribute(resources.AttName)
		t0 := task.SubM(reporter.NewMessage("trigger for s3 bucket").
			WithArg("bucketName", bucketName).
			WithArg("bucketArn", s3Resource.GetAttribute(resources.AttArn)))
		var notification = &awsapi.S3NotificationDefinition{
			Lambdas: make([]awsapi.S3NotificationLambdaDefinition, 0, 4),
		}
		for _, lambdaTrigger := range s3CopyTrigger.Lambdas {
			lambdaRes := nestorResources.FindResourceByID(lambdaTrigger.LambdaID)
			lambdaArn := lambdaRes.GetAttribute(resources.AttArn)
			if lambdaRes == nil {
				return errors.New("s3uploadTrigger: lambda not registered:" + lambdaTrigger.LambdaID)
			}
			err = api.GiveS3LambdaInvokePermission(lambdaArn, bucketArn, bucketName, t0)
			if err != nil {
				return err
			}

			def := awsapi.S3NotificationLambdaDefinition{
				LambdaArn: lambdaRes.GetAttribute(resources.AttArn),
			}
			if lambdaTrigger.Prefix != "" {
				def.Prefix = lambdaTrigger.Prefix
			}
			if lambdaTrigger.Suffix != "" {
				def.Suffix = lambdaTrigger.Suffix
			}
			notification.Lambdas = append(notification.Lambdas, def)
		}
		err = api.SetBucketNotificationConfiguration(bucketName, notification, t0)
		if err != nil {
			return err
		}
	}
	return nil
}
