package actions

import (
	"errors"
	"fmt"

	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
)

func (actions *Actions) createS3UploadTrigger(s3Trigger config.TriggerS3CopyDefinition) error {
	var nestorResources = actions.nestorResources
	s3UploadRes := nestorResources.FindResourceByID(s3Trigger.BucketID)
	if s3UploadRes == nil {
		return errors.New("s3uploadTrigger: bucket not registered:" + s3Trigger.BucketID)
	}
	lambdaRes := nestorResources.FindResourceByID(s3Trigger.LambdaID)
	if lambdaRes == nil {
		return errors.New("s3uploadTrigger: lambda not registered:" + s3Trigger.LambdaID)
	}
	fmt.Printf("ResS3BucketForUpload: %#v", s3UploadRes)
	fmt.Printf("lambdaRes: %#v", lambdaRes)
	return nil
}

// CreateTriggers processing for provision CLI command
func (actions *Actions) CreateTriggers(task *reporter.Task) error {
	// var appName = actions.nestorConfig.Application.Name
	// var environment = actions.environment
	// var api = actions.api
	var nestorConfig = actions.nestorConfig
	// var nestorResources = actions.nestorResources

	var err error
	var triggers = nestorConfig.Triggers
	// s3 triggers
	for _, s3CopyTrigger := range triggers.S3copy {
		err = actions.createS3UploadTrigger(s3CopyTrigger)
		if err != nil {
			return err
		}
	}
	return nil
}
