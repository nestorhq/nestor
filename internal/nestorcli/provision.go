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
	_, err1 := api.CreateUserPool(userPoolName, "nestor.res.cognito.main", t1)
	if err1 != nil {
		t1.Fail(err1)
		panic(err1)
	}
	t1.Ok()

	// 2: dynamodb table
	var tableName = appName + "-" + environment + "-main"
	t2 := t.SubM(reporter.NewMessage("create dynamodb table").WithArg("tableName", tableName))
	_, err2 := api.CreateMonoTable(tableName, "nestor.res.dynamoDbTable.main", t2)
	if err2 != nil {
		t2.Fail(err2)
		panic(err2)
	}
	t2.Ok()

	// 3: event bus
	var eventBusName = appName + "-" + environment + "-main"
	t3 := t.SubM(reporter.NewMessage("create event bus").WithArg("eventBusName", eventBusName))
	_, err3 := api.CreateEventBus(eventBusName, "nestor.res.eventBridge.main", t3)
	if err3 != nil {
		t3.Fail(err3)
		panic(err3)
	}
	t3.Ok()

	// 4: s3 bucket - storage
	var bucketNameStorage = appName + "-" + environment + "-store"
	t4 := t.SubM(reporter.NewMessage("create s3 bucket for storage").WithArg("bucketNameStorage", bucketNameStorage))
	_, err4 := api.CreateBucket(bucketNameStorage, "nestor.res.s3.storage", t4)
	if err4 != nil {
		t4.Fail(err4)
		panic(err4)
	}
	t4.Ok()

	// 5: s3 bucket - upload
	var bucketNameUpload = appName + "-" + environment + "-upload"
	t5 := t.SubM(reporter.NewMessage("create s3 bucket for upload").WithArg("bucketNameUpload", bucketNameUpload))
	_, err5 := api.CreateBucket(bucketNameUpload, "nestor.res.s3.upload", t5)
	if err5 != nil {
		t5.Fail(err5)
		panic(err5)
	}
	t5.Ok()

	t.Ok()
}
