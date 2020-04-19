package actions

import (
	"errors"

	"github.com/nestorhq/nestor/internal/reporter"
	"github.com/nestorhq/nestor/internal/resources"
)

// CreateResources processing for provision CLI command
func (actions *Actions) CreateResources(task *reporter.Task) error {
	var appName = actions.nestorConfig.Application.Name
	var environment = actions.environment
	var api = actions.api
	var nestorConfig = actions.nestorConfig
	var nestorResources = actions.nestorResources

	var t = task.SubM(reporter.NewMessage("CreateResources").
		WithArg("environment", environment).
		WithArg("appName", appName))

	t.Section("Create Cognito User Pools")
	for _, userPool := range nestorConfig.Resources.CognitoUserPool {
		var userPoolName = appName + "-" + environment + "-" + userPool.ID
		var resourceID = "resources.cognito_userpool." + userPool.ID
		t1 := t.SubM(reporter.NewMessage("create user pool").
			WithArg("userPoolName", userPoolName).
			WithArg("resourceID", resourceID))
		res, err := api.CreateUserPool(userPoolName, resourceID, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		nestorResources.RegisterNestorResource(resourceID, resources.CognitoUserPool, resources.AttArn, res.UserPoolArn)
		t1.Ok()
	}

	t.Section("Create DynamoDb tables")
	for _, dynamoDbTable := range nestorConfig.Resources.DynamodbTable {
		var tableName = appName + "-" + environment + "-" + dynamoDbTable.ID
		var resourceID = "resources.dynamodb_table." + dynamoDbTable.ID
		t1 := t.SubM(reporter.NewMessage("create dynamodb table").
			WithArg("tableName", tableName).
			WithArg("resourceID", resourceID))
		res, err := api.CreateMonoTable(tableName, resourceID, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		nestorResources.RegisterNestorResource(resourceID, resources.DynamoDbTable, resources.AttArn, res.TableArn)
		t1.Ok()
	}

	t.Section("Create EventBridge buses")
	for _, eventBridgeBus := range nestorConfig.Resources.EventBridgeBus {
		var eventBridgeBusName = appName + "-" + environment + "-" + eventBridgeBus.ID
		var resourceID = "resources.eventbridge_bus." + eventBridgeBus.ID
		t1 := t.SubM(reporter.NewMessage("create eventBridge bus").
			WithArg("eventBridgeBusName", eventBridgeBusName).
			WithArg("resourceID", resourceID))
		res, err := api.CreateEventBus(eventBridgeBusName, resourceID, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		nestorResources.RegisterNestorResource(resourceID, resources.EventBridgeBus, resources.AttArn, res.EventBusArn)
		t1.Ok()
	}

	t.Section("Create S3 buckets")
	for _, s3Bucket := range nestorConfig.Resources.S3Bucket {
		var s3BucketName = appName + "-" + environment + "-" + s3Bucket.ID
		var resourceID = "resources.s3_bucket." + s3Bucket.ID
		t1 := t.SubM(reporter.NewMessage("create s3 bucket").
			WithArg("s3BucketName", s3BucketName).
			WithArg("resourceID", resourceID))
		res, err := api.CreateBucket(s3BucketName, resourceID, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		nestorResources.RegisterNestorResource(resourceID, resources.S3Bucket, resources.AttArn, res.BucketArn)
		nestorResources.RegisterNestorResource(resourceID, resources.S3Bucket, resources.AttName, res.BucketName)
		t1.Ok()
	}

	t.Section("Create CloudwatchLogs groups")
	for _, cloudwatchLogGroup := range nestorConfig.Resources.CloudwatchlogsGroup {
		var cloudwatchLogGroupName = appName + "-" + environment + "-" + cloudwatchLogGroup.ID
		var resourceID = "resources.cloudwatchlogs_group." + cloudwatchLogGroup.ID
		t1 := t.SubM(reporter.NewMessage("create cloudwatLogs group").
			WithArg("cloudwatchLogGroupName", cloudwatchLogGroupName).
			WithArg("resourceID", resourceID))
		res, err := api.CreateCloudWatchGroup(cloudwatchLogGroupName, resourceID, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		nestorResources.RegisterNestorResource(resourceID, resources.CloudwatchLogGroup, resources.AttName, res.GroupName)
		t1.Ok()
	}

	t.Section("Create Lambda Functions")
	for _, lambdaFunction := range nestorConfig.Resources.LambdaFunction {
		var lambdaFunctionName = appName + "-" + environment + "-" + lambdaFunction.ID
		var roleName = appName + "-" + environment + "-" + lambdaFunction.ID
		var resourceID = "resources.lambda_function." + lambdaFunction.ID
		t1 := t.SubM(reporter.NewMessage("create lambda role").
			WithArg("lambdaFunctionName", lambdaFunctionName).
			WithArg("roleName", roleName).
			WithArg("resourceID", resourceID))

		role, err := api.CreateAppLambdaRole(roleName, resourceID, lambdaFunctionName, lambdaFunction, actions.nestorResources, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		t1.LogM(reporter.NewMessage("role created").
			WithArg("RoleArn", role.RoleArn).WithArg("RoleName", role.RoleName))
		t1.Ok()

		t2 := t.SubM(reporter.NewMessage("create Lambda ").WithArg("lambdaFunctionName", lambdaFunctionName))
		res, err2 := api.CreateLambda(lambdaFunctionName, resourceID, role.RoleArn, lambdaFunction.Runtime, t1)
		if err2 != nil {
			t2.Fail(err2)
			return err2
		}
		nestorResources.RegisterNestorResource(resourceID, resources.LambdaFunction, resources.AttArn, res.FunctionArn)
		nestorResources.RegisterNestorResource(resourceID, resources.LambdaFunction, resources.AttName, res.FunctionName)
		t2.Ok()
	}

	t.Section("Create ApiGateway HTTP")
	for _, apiGatewayHTTP := range nestorConfig.Resources.ApigatewayHTTP {
		var apiGatewayHTTPName = appName + "-" + environment + "-" + apiGatewayHTTP.ID
		var lambdaTarget = nestorResources.FindResourceByID(apiGatewayHTTP.TargetLambdaID)
		var resourceID = "resources.apigateway_http." + apiGatewayHTTP.ID
		t1 := t.SubM(reporter.NewMessage("create rest API").
			WithArg("apiGatewayHTTPName", apiGatewayHTTPName).
			WithArg("resourceID", resourceID))
		if lambdaTarget == nil {
			err := errors.New("can't find target for api gateway:" + apiGatewayHTTP.TargetLambdaID)
			t1.Fail(err)
			return err
		}
		lambdaTargetArn := lambdaTarget.GetAttribute(resources.AttArn)
		res, err := api.CreateRestAPI(apiGatewayHTTPName, lambdaTargetArn, resourceID, t1)
		if err != nil {
			t1.Fail(err)
			return err
		}
		nestorResources.RegisterNestorResource(resourceID, resources.HTTPAPIGateway, resources.AttID, res.HTTPApiID)
		t1.Ok()

		t2 := t.SubM(reporter.NewMessage("give Api invoke permission").
			WithArg("lambdaTargetArn", lambdaTargetArn).
			WithArg("apiID", res.HTTPApiID))
		api.GiveAPIGatewayLambdaInvokePermission(lambdaTargetArn, res.HTTPApiID, t2)
		t2.Ok()
	}

	t.Ok()
	return nil
}
