package nestorcli

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
)

// CliProvision processing for provision CLI command
func CliProvision(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.App.Name

	var t = reporter.NewReporterM(reporter.NewMessage("command: provision").
		WithArg("environment", environment).
		WithArg("appName", appName).
		WithArg("config", fmt.Sprintf("%v", nestorConfig))).
		Start()

	// TODO: hard coded
	resourceTags := awsapi.NewResourceTag("1", environment, appName)

	// TODO: hard coded
	t0 := t.Sub("Initialize aws API")
	api, err := awsapi.NewAwsAPI("sls", resourceTags, "us-west-1", "us-west-2", t0)

	if err != nil {
		panic(err)
	}

	// 1: user pool
	var userPoolName = appName + "-" + environment
	t1 := t.SubM(reporter.NewMessage("create user pool").WithArg("userPoolName", userPoolName))
	_, errup := api.CreateUserPool(userPoolName, t1)
	if errup != nil {
		t1.Fail(errup)
		panic(errup)
	}
	t1.Ok()

	// 2: dynamodb tables
	var tableName = appName + "-" + environment + "-main"
	t2 := t.SubM(reporter.NewMessage("create dynamodb table").WithArg("tableName", tableName))
	_, errdynamo := api.CreateMonoTable(tableName, t2)
	if errdynamo != nil {
		t2.Fail(errdynamo)
		panic(errdynamo)
	}
	t2.Ok()

	t.Ok()
}
