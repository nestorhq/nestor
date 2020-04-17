package nestorcli

import (
	"fmt"

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
		WithArg("appName", appName).
		WithArg("config", fmt.Sprintf("%#v", nestorConfig))).
		Start()

	// TODO: hard coded
	resourceTags := awsapi.NewResourceTag("1", environment, appName)

	// TODO: hard coded
	t0 := t.Sub("Initialize aws API")
	api, err := awsapi.NewAwsAPI("sls", resourceTags, "us-west-1", "us-west-2", t0)

	if err != nil {
		panic(err)
	}
	err = actions.CreateResources(environment, api, nestorConfig, nestorResources, t)
	if err != nil {
		panic(err)
	}

}
