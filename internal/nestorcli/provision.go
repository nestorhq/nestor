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

	// user pool
	if ok, resID := nestorConfig.IsResourceRequired(config.ResCognitoMain); ok {
		t.Section("configure resource:" + resID)
		var userPoolName = appName + "-" + environment
		t1 := t.SubM(reporter.NewMessage("create user pool").WithArg("userPoolName", userPoolName))
		res, err := api.CreateUserPool(userPoolName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.UserPoolArn)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// dynamodb table
	if ok, resID := nestorConfig.IsResourceRequired(config.ResDynamoDbTableMain); ok {
		t.Section("configure resource:" + resID)
		var tableName = appName + "-" + environment + "-main"
		t1 := t.SubM(reporter.NewMessage("create dynamodb table").WithArg("tableName", tableName))
		res, err := api.CreateMonoTable(tableName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.TableArn)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// event bus
	if ok, resID := nestorConfig.IsResourceRequired(config.ResEventBridgeMain); ok {
		t.Section("configure resource:" + resID)
		var eventBusName = appName + "-" + environment + "-main"
		t1 := t.SubM(reporter.NewMessage("create event bus").WithArg("eventBusName", eventBusName))
		res, err := api.CreateEventBus(eventBusName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.EventBusArn)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// s3 bucket - storage
	if ok, resID := nestorConfig.IsResourceRequired(config.ResS3BucketForStorage); ok {
		t.Section("configure resource:" + resID)
		var bucketNameStorage = appName + "-" + environment + "-store"
		t1 := t.SubM(reporter.NewMessage("create s3 bucket for storage").WithArg("bucketNameStorage", bucketNameStorage))
		res, err := api.CreateBucket(bucketNameStorage, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.BucketArn)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// s3 bucket - upload
	if ok, resID := nestorConfig.IsResourceRequired(config.ResS3BucketForUpload); ok {
		t.Section("configure resource:" + resID)
		var bucketNameUpload = appName + "-" + environment + "-upload"
		t1 := t.SubM(reporter.NewMessage("create s3 bucket for upload").WithArg("bucketNameUpload", bucketNameUpload))
		res, err := api.CreateBucket(bucketNameUpload, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.BucketArn)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}
	// CreateRestAPI
	if ok, resID := nestorConfig.IsResourceRequired(config.ResHTTPAPIMain); ok {
		t.Section("configure resource:" + resID)
		var restAPIName = appName + "-" + environment + "-main"
		t1 := t.SubM(reporter.NewMessage("create Rest API").WithArg("restAPIName", restAPIName))
		res, err := api.CreateRestAPI(restAPIName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.HTTPApiID)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}
	// CreateLogGroup
	if ok, resID := nestorConfig.IsResourceRequired(config.ResLogGroupMainEventBridge); ok {
		t.Section("configure resource:" + resID)
		var groupName = "/aws/events/" + appName + "/" + environment + "/mainEventBridgeTarget"
		t1 := t.SubM(reporter.NewMessage("create CloudWatchGroup ").WithArg("groupName", groupName))
		res, err := api.CreateCloudWatchGroup(groupName, resID, t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		nestorConfig.RegisterNestorResource(resID, res.GroupName)
		t1.Ok()
	} else {
		t0.Log("non required resource:" + resID)
	}

	// Create lambdas
	for _, lambda := range nestorConfig.Lambdas {
		var lambdaName = appName + "-" + environment + "-" + lambda.ID
		t.Section("create lambda:" + lambdaName)
		t1 := t.SubM(reporter.NewMessage("create Lambda ").WithArg("lambdaName", lambdaName))
		_, err := api.CreateLambda(lambdaName, lambda.ID, "arn:aws:iam::464972470401:role/iam_for_events_lambda_prod", t1)
		if err != nil {
			t1.Fail(err)
			panic(err)
		}
		t1.Ok()
	}

	t.Ok()
}
