package nestorcli

import (
	"github.com/nestorhq/nestor/internal/actions"
	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
	"github.com/nestorhq/nestor/internal/resources"
)

// CliProvision processing for provision CLI command
func CliProvision(environment string, nestorConfig *config.Config) {
	var nestorResources = resources.NewResources()
	var appName = nestorConfig.Application.Name

	var t = reporter.NewReporterM(reporter.NewMessage("command: provision").
		WithArg("environment", environment).
		WithArg("appName", appName)).
		Start()

	// TODO: hard coded
	resourceTags := awsapi.NewResourceTag("1", environment, appName)

	// TODO: hard coded
	t0 := t.Sub("Initialize aws API")
	api, err := awsapi.NewAwsAPI(nestorConfig.Application.ProfileName,
		resourceTags, nestorConfig.Application.Region, nestorConfig.Application.RegionCognito, t0)

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
