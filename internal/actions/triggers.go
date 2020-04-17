package actions

import (
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
)

func (actions *Actions) createS3UploadTrigger(s3Trigger config.TriggerS3UploadDefinition) error {
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
	// s3 triggers
	for _, trigger := range nestorConfig.Triggers {
		for _, s3UploadTrigger := range trigger.S3upload {
			err = actions.createS3UploadTrigger(s3UploadTrigger)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
