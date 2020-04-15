package nestorcli

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
)

// CliProvision processing for provision CLI command
func CliProvision(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.Application.Name

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
	if ok, resID := nestorConfig.IsResourceRequired(config.ResCognitoMain); ok {
		t0.Section("configure resource:" + resID)
		var userPoolName = appName + "-" + environment
		t1 := t.SubM(reporter.NewMessage("create user pool").WithArg("userPoolName", userPoolName))
		_, err := api.CreateUserPool(userPoolName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// 2: dynamodb table
	if ok, resID := nestorConfig.IsResourceRequired(config.ResDynamoDbTableMain); ok {
		t0.Section("configure resource:" + resID)
		var tableName = appName + "-" + environment + "-main"
		t1 := t.SubM(reporter.NewMessage("create dynamodb table").WithArg("tableName", tableName))
		_, err := api.CreateMonoTable(tableName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// 3: event bus
	if ok, resID := nestorConfig.IsResourceRequired(config.ResEventBridgeMain); ok {
		t0.Section("configure resource:" + resID)
		var eventBusName = appName + "-" + environment + "-main"
		t1 := t.SubM(reporter.NewMessage("create event bus").WithArg("eventBusName", eventBusName))
		_, err := api.CreateEventBus(eventBusName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// 4: s3 bucket - storage
	if ok, resID := nestorConfig.IsResourceRequired(config.ResS3BucketForStorage); ok {
		t0.Section("configure resource:" + resID)
		var bucketNameStorage = appName + "-" + environment + "-store"
		t1 := t.SubM(reporter.NewMessage("create s3 bucket for storage").WithArg("bucketNameStorage", bucketNameStorage))
		_, err := api.CreateBucket(bucketNameStorage, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// 5: s3 bucket - upload
	if ok, resID := nestorConfig.IsResourceRequired(config.ResS3BucketForUpload); ok {
		t0.Section("configure resource:" + resID)
		var bucketNameUpload = appName + "-" + environment + "-upload"
		t1 := t.SubM(reporter.NewMessage("create s3 bucket for upload").WithArg("bucketNameUpload", bucketNameUpload))
		_, err := api.CreateBucket(bucketNameUpload, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}
	//6: CreateRestAPI
	if ok, resID := nestorConfig.IsResourceRequired(config.ResHTTPAPIMain); ok {
		t0.Section("configure resource:" + resID)
		var restAPIName = appName + "-" + environment + "-main"
		t1 := t.SubM(reporter.NewMessage("create Rest API").WithArg("restAPIName", restAPIName))
		_, err := api.CreateRestAPI(restAPIName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}
	//7: CreateLogGroup
	if ok, resID := nestorConfig.IsResourceRequired(config.ResLogGroupMainEventBridge); ok {
		t0.Section("configure resource:" + resID)
		var groupName = appName + "-" + environment + "-mainEventBridgeTarget"
		t1 := t.SubM(reporter.NewMessage("create CloudWatchGroup ").WithArg("groupName", groupName))
		_, err := api.CreateCloudWatchGroup(groupName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}
	t.Ok()
}
