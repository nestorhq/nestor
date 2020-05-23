package awsapi

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/nestorhq/nestor/internal/reporter"
	_ "github.com/nestorhq/nestor/internal/templates/statik" // embedded fs
	"github.com/rakyll/statik/fs"
)

// LambdaAPI api
type LambdaAPI struct {
	resourceTags *ResourceTags
	client       *lambda.Lambda
	account      string
}

// LambdaInformation description of a lambda
type LambdaInformation struct {
	FunctionName string
	FunctionArn  string
}

// LambdaCreateInformation description of create information for lambda
type LambdaCreateInformation struct {
	runtime     string
	templateZip string
	handler     string
}

// NewLambdaAPI constructor
func NewLambdaAPI(session *session.Session, resourceTags *ResourceTags, account string) (*LambdaAPI, error) {
	var api = LambdaAPI{resourceTags: resourceTags, account: account}
	// Create Lambda client
	api.client = lambda.New(session)
	return &api, nil
}

func (api *LambdaAPI) updateLambdaConfiguration(lambdaName string, handler string,
	environmentVariables map[string]string, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.UpdateFunctionConfiguration").WithArg("lambdaName", lambdaName))

	input := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(lambdaName),
		Handler:      aws.String(handler),
		Environment: &lambda.Environment{
			Variables: aws.StringMap(environmentVariables),
		},
	}
	result, err := api.client.UpdateFunctionConfiguration(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	t0.Okr(map[string]string{
		"FunctionArn":  *result.FunctionArn,
		"FunctionName": *result.FunctionName,
		"Version":      *result.Version,
		"State":        *result.State,
	})

	return &LambdaInformation{
		FunctionName: *result.FunctionName,
		FunctionArn:  *result.FunctionArn,
	}, nil

}

func (api *LambdaAPI) updateLambdaCode(lambdaName string, zipFileName string, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.UpdateFunctionCode").WithArg("lambdaName", lambdaName))

	zipData, err := ioutil.ReadFile(zipFileName)
	if err != nil {
		return nil, err
	}
	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(lambdaName),
		Publish:      aws.Bool(true),
		ZipFile:      zipData,
	}
	result, err := api.client.UpdateFunctionCode(input)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	t0.Okr(map[string]string{
		"FunctionArn":  *result.FunctionArn,
		"FunctionName": *result.FunctionName,
		"Version":      *result.Version,
		"State":        *result.State,
	})

	return &LambdaInformation{
		FunctionName: *result.FunctionName,
		FunctionArn:  *result.FunctionArn,
	}, nil

}

func (api *LambdaAPI) doCreateLambda(lambdaName string, nestorID string, roleArn string, createInformation *LambdaCreateInformation, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.CreateFunction").WithArg("lambdaName", lambdaName))

	// read lambda code from statik (embedded) filesystem
	t1 := t0.SubM(reporter.NewMessage("create lambda with statik").
		WithArg("lambdaName", createInformation.runtime).
		WithArg("templateZip", createInformation.templateZip).
		WithArg("handler", createInformation.handler))
	statikFS, err := fs.New()
	if err != nil {
		t1.Log("cannot initialize statikFS")
		t1.Fail(err)
		return nil, err
	}
	r, err := statikFS.Open(createInformation.templateZip)
	// fmt.Printf("@@ statik file: %#v\n", r)
	if err != nil {
		t1.Log("cannot read from statikFS:" + createInformation.templateZip)
		t1.Fail(err)
		return nil, err
	}
	defer r.Close()
	zipData, err := ioutil.ReadAll(r)
	if err != nil {
		t1.Log("cannot read all data from statikFS:" + createInformation.templateZip)
		t1.Fail(err)
		return nil, err
	}

	createCode := &lambda.FunctionCode{
		ZipFile: zipData,
	}

	input := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String(lambdaName),
		Tags:         aws.StringMap(api.resourceTags.getTagsAsMapWithID(nestorID)),
		Handler:      aws.String(createInformation.handler),
		Runtime:      aws.String(createInformation.runtime),
		Role:         aws.String(roleArn),
	}
	result, err := api.client.CreateFunction(input)
	if err != nil {
		// we don't log error for that AWS error as we may be in a retry loop
		if getAwsErrorCode(err) == "InvalidParameterValueException" {
			return nil, err
		}
		t0.Fail(err)
		return nil, err
	}
	t0.Okr(map[string]string{
		"FunctionArn":  *result.FunctionArn,
		"FunctionName": *result.FunctionName,
		"Version":      *result.Version,
		"State":        *result.State,
	})

	return &LambdaInformation{
		FunctionName: *result.FunctionName,
		FunctionArn:  *result.FunctionArn,
	}, nil

}

func (api *LambdaAPI) checkLambdaExistence(lambdaName string, task *reporter.Task) (*LambdaInformation, map[string]*string, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.GetFunction").WithArg("lambdaName", lambdaName))
	input := &lambda.GetFunctionInput{
		FunctionName: aws.String(lambdaName),
	}
	result, err := api.client.GetFunction(input)
	if err != nil {
		if getAwsErrorCode(err) == "ResourceNotFoundException" {
			return nil, nil, nil
		}
		t0.Fail(err)
		return nil, nil, err
	}
	t0.Okr(map[string]string{
		"FunctionArn":  *result.Configuration.FunctionArn,
		"FunctionName": *result.Configuration.FunctionName,
		"Version":      *result.Configuration.Version,
		"State":        *result.Configuration.State,
	})

	return &LambdaInformation{
		FunctionName: *result.Configuration.FunctionName,
		FunctionArn:  *result.Configuration.FunctionArn,
	}, result.Tags, nil

}

func (api *LambdaAPI) checkLambdaExistenceAndTags(lambdaName string, nestorID string, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("checkLambdaExistenceAndTags").WithArg("lambdaName", lambdaName))
	lambdaInformation, tags, err := api.checkLambdaExistence(lambdaName, t0)
	if err != nil {
		t0.Fail(err)
		return nil, err
	}
	if lambdaInformation == nil {
		t0.Ok()
		return nil, nil
	}

	t1 := task.SubM(reporter.NewMessage("checkLambdaTags").WithArg("lambdaName", lambdaName))
	err2 := api.resourceTags.checkTags(tags, nestorID)
	if err2 != nil {
		t1.Fail(err2)
		return nil, err2
	}
	return lambdaInformation, nil
}

func (api *LambdaAPI) createLambda(lambdaName string, nestorID string, roleArn string, createInformation *LambdaCreateInformation, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("createLambda").WithArg("lambdaName", lambdaName))

	t1 := t0.Sub("check if lambda exists")
	lambdaInformation, err := api.checkLambdaExistenceAndTags(lambdaName, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}

	if lambdaInformation != nil {
		t1.Log("lambda exists")
		t1.Okr(map[string]string{
			"functionName": lambdaInformation.FunctionName,
			"arn":          lambdaInformation.FunctionArn,
		})

		return lambdaInformation, nil
	}

	t2 := t0.Sub("lambda does not exist - creating it")
	result, err := api.doCreateLambda(lambdaName, nestorID, roleArn, createInformation, t2)
	if err != nil {
		t2.Fail(err)
		return nil, err
	}
	t2.Ok()
	t0.Okr(map[string]string{
		"functionName": result.FunctionName,
		"arn":          result.FunctionArn,
	})
	return result, nil
}

func (api *LambdaAPI) addInvokePermission(lambdaArn string, sid string, principal string, sourceArn string, sourceAccount string, task *reporter.Task) error {
	t0 := task.SubM(reporter.NewMessage("addInvokePermission").
		WithArg("lambdaArn", lambdaArn).
		WithArg("sid", sid))

	input := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(lambdaArn),
		Principal:    aws.String(principal),
		StatementId:  aws.String(sid),
		SourceArn:    aws.String(sourceArn),
	}
	if sourceAccount != "" {
		input.SourceAccount = aws.String(sourceAccount)
	}
	_, err := api.client.AddPermission(input)
	if err != nil {
		// if getAwsErrorCode(err) == "ResourceNotFoundException" {
		// 	return nil
		// }
		t0.Fail(err)
		return err
	}
	t0.Ok()
	return nil
}

func (api *LambdaAPI) removePermission(lambdaArn string, sid string, task *reporter.Task) error {
	t0 := task.SubM(reporter.NewMessage("removePermission").
		WithArg("lambdaArn", lambdaArn).
		WithArg("sid", sid))

	input := &lambda.RemovePermissionInput{
		FunctionName: aws.String(lambdaArn),
		StatementId:  aws.String(sid),
	}
	_, err := api.client.RemovePermission(input)
	if err != nil {
		if getAwsErrorCode(err) == "ResourceNotFoundException" {
			t0.Ok()
			return nil
		}
		t0.Fail(err)
		return err
	}
	t0.Ok()
	return nil
}

func (api *LambdaAPI) giveS3InvokePermission(lambdaArn string, bucketArn string, bucketName string, task *reporter.Task) error {
	var sid = "sid-s3invoke-" + bucketName
	t0 := task.SubM(reporter.NewMessage("giveS3InvokePermission").
		WithArg("lambdaArn", lambdaArn).
		WithArg("bucketArn", bucketArn))
	err := api.removePermission(lambdaArn, sid, t0)
	if err != nil {
		t0.Fail(err)
		return err
	}

	err = api.addInvokePermission(lambdaArn, sid, "s3.amazonaws.com", bucketArn, api.account, t0)
	if err != nil {
		t0.Fail(err)
		return err
	}
	t0.Ok()
	return nil
}

func (api *LambdaAPI) giveAPIGatewayInvokePermission(lambdaArn string, apiID string, task *reporter.Task) error {
	var sid = "sid-apigwinvoke-" + apiID
	t0 := task.SubM(reporter.NewMessage("giveApiGatewayInvokePermission").
		WithArg("lambdaArn", lambdaArn).
		WithArg("apiID", apiID))
	err := api.removePermission(lambdaArn, sid, t0)
	if err != nil {
		t0.Fail(err)
		return err
	}

	sourceArn := fmt.Sprintf("arn:aws:execute-api:us-west-2:%s:%s/*/$default", api.account, apiID)

	err = api.addInvokePermission(lambdaArn, sid, "apigateway.amazonaws.com", sourceArn, "", t0)
	if err != nil {
		t0.Fail(err)
		return err
	}
	t0.Ok()
	return nil
}
