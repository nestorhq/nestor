package awsapi

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/nestorhq/nestor/internal/reporter"
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

const defaultLambda = `exports.handler = async (event, context, callback) => {
  console.log('>lambda>event> ', JSON.stringify(event, null, '  '));
  callback();
}`

func makeZipData(content string, fileName string) ([]byte, error) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{fileName, content},
	}
	for _, file := range files {
		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = zipFile.Write([]byte(file.Body))
		if err != nil {
			return nil, err
		}
	}

	// Make sure to check the error on Close.
	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// NewLambdaAPI constructor
func NewLambdaAPI(session *session.Session, resourceTags *ResourceTags, account string) (*LambdaAPI, error) {
	var api = LambdaAPI{resourceTags: resourceTags, account: account}
	// Create Lambda client
	api.client = lambda.New(session)
	return &api, nil
}

func (api *LambdaAPI) doCreateLambda(lambdaName string, nestorID string, roleArn string, runtime string, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("api.client.CreateFunction").WithArg("lambdaName", lambdaName))

	zipData, err := makeZipData(defaultLambda, "index.js")
	if err != nil {
		return nil, err
	}

	createCode := &lambda.FunctionCode{
		ZipFile: zipData,
	}

	input := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String(lambdaName),
		Tags:         aws.StringMap(api.resourceTags.getTagsAsMapWithID(nestorID)),
		Handler:      aws.String("index.handler"),
		Runtime:      aws.String(runtime),
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

func (api *LambdaAPI) createLambda(lambdaName string, nestorID string, roleArn string, runtime string, task *reporter.Task) (*LambdaInformation, error) {
	t0 := task.SubM(reporter.NewMessage("createLambda").WithArg("lambdaName", lambdaName))

	t1 := t0.Sub("check if lambda exists")
	lambdaInformation, err := api.checkLambdaExistenceAndTags(lambdaName, nestorID, t1)
	if err != nil {
		t1.Fail(err)
		return nil, err
	}

	if lambdaInformation != nil {
		t1.Log("table exists")
		t1.Okr(map[string]string{
			"functionName": lambdaInformation.FunctionName,
			"arn":          lambdaInformation.FunctionArn,
		})

		return lambdaInformation, nil
	}

	t2 := t0.Sub("lambda does not exist - creating it")
	result, err := api.doCreateLambda(lambdaName, nestorID, roleArn, runtime, t2)
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
	t0 := task.SubM(reporter.NewMessage("giveS3InvokePermission").
		WithArg("lambdaArn", lambdaArn).
		WithArg("apiID", apiID))
	err := api.removePermission(lambdaArn, sid, t0)
	if err != nil {
		t0.Fail(err)
		return err
	}

	sourceArn := fmt.Sprintf("arn:aws:execute-api:us-west-1:%s:%s/*/$default", api.account, apiID)

	err = api.addInvokePermission(lambdaArn, sid, "apigateway.amazonaws.com", sourceArn, "", t0)
	if err != nil {
		t0.Fail(err)
		return err
	}
	t0.Ok()
	return nil
}
