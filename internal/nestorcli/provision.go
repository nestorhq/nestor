package nestorcli

import (
	"github.com/nestorhq/nestor/internal/actions"
	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
	"github.com/nestorhq/nestor/internal/resources"
)

// NestorConfigVersion version of the nestor config format
const NestorConfigVersion = "1"

// CliProvision processing for provision CLI command
func CliProvision(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.Application.Name
	var profileName = nestorConfig.Application.ProfileName
	var region = nestorConfig.Application.Region
	var regionCognito = nestorConfig.Application.RegionCognito
	var nestorResources = resources.NewResources()

	var t = reporter.NewReporterM(reporter.NewMessage("command: provision").
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
	awsActions := actions.NewActions(environment, api, nestorConfig, nestorResources)

	t1 := t.Sub("create resources")
	err = awsActions.CreateResources(t1)
	if err != nil {
		t1.Fail(err)
		panic(err)
	}
	t1.Ok()

	t2 := t.Sub("create triggers")
	err = awsActions.CreateTriggers(t2)
	if err != nil {
		t2.Fail(err)
		panic(err)
	}

}
