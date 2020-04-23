package nestorcli

import (
	"github.com/nestorhq/nestor/internal/actions"
	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
	"github.com/nestorhq/nestor/internal/resources"
)

// CliDeploy processing for deploy CLI command
func CliDeploy(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.Application.Name
	var profileName = nestorConfig.Application.ProfileName
	var region = nestorConfig.Application.Region
	var regionCognito = nestorConfig.Application.RegionCognito

	var t = reporter.NewReporterM(reporter.NewMessage("command: deploy").
		WithArg("environment", environment).
		WithArg("appName", appName)).
		Start()
	// tags to associate to the created aws resources
	resourceTags := awsapi.NewResourceTag(NestorConfigVersion, environment, appName)

	t0 := t.Sub("Initialize aws API")
	api, err := awsapi.NewAwsAPI(profileName, resourceTags, region, regionCognito, t0)

	if err != nil {
		panic(err)
	}

	var nestorResources = resources.NewResources()
	awsActions := actions.NewActions(environment, api, nestorConfig, nestorResources)

	t1 := t.Sub("create resources")
	err = awsActions.DoDeployment(t1)
	if err != nil {
		t1.Fail(err)
		panic(err)
	}
	t1.Ok()

	t.Ok()
}
