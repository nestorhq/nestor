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
		return errors.New("ResS3BucketForUpload resource not registered")
	}
	fmt.Printf("ResS3BucketForUpload: %#v", s3UploadRes)
	return nil
}

// CreateTriggers processing for provision CLI command
func (actions *Actions) CreateTriggers(task *reporter.Task) error {
	// var appName = actions.nestorConfig.Application.Name
	// var environment = actions.environment
	// var api = actions.api
	//	var nestorConfig = actions.nestorConfig
	// var nestorResources = actions.nestorResources

	// var err error
	// // s3 triggers
	// for _, trigger := range nestorConfig.Triggers {
	// 	for _, s3CopyTrigger := range trigger.S3copy {
	// 		err = actions.createS3UploadTrigger(s3CopyTrigger)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	return nil
}
