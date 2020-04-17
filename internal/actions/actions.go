package actions

import (
	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/resources"
)

// Actions actions holder
type Actions struct {
	environment     string
	api             *awsapi.AwsAPI
	nestorConfig    *config.Config
	nestorResources *resources.Resources
}

// NewActions ctor
func NewActions(environment string, api *awsapi.AwsAPI, nestorConfig *config.Config, nestorResources *resources.Resources) *Actions {
	return &Actions{
		environment:     environment,
		api:             api,
		nestorConfig:    nestorConfig,
		nestorResources: nestorResources,
	}
}
