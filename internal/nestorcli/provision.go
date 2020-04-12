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
	_, err1 := api.CreateUserPool(userPoolName, t1)
	if err1 != nil {
		t1.Fail(err1)
		panic(err1)
	}
	t1.Ok()

	// 2: dynamodb tables
	var tableName = appName + "-" + environment + "-main"
	t2 := t.SubM(reporter.NewMessage("create dynamodb table").WithArg("tableName", tableName))
	_, err2 := api.CreateMonoTable(tableName, t2)
	if err2 != nil {
		t2.Fail(err2)
		panic(err2)
	}
	t2.Ok()

	// 3: event bus
	var eventBusName = appName + "-" + environment + "-main"
	t3 := t.SubM(reporter.NewMessage("create event bus").WithArg("eventBusName", eventBusName))
	_, err3 := api.CreateEventBus(eventBusName, t3)
	if err3 != nil {
		t3.Fail(err3)
		panic(err3)
	}
	t3.Ok()

	t.Ok()
}
