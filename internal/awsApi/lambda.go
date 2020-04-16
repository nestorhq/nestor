package awsapi

import (
	"archive/zip"
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/nestorhq/nestor/internal/reporter"
)

// LambdaAPI api
type LambdaAPI struct {
	resourceTags *ResourceTags
	client       *lambda.Lambda
}

// LambdaInformation description of a lambda
type LambdaInformation struct {
	functionName string
	arn          string
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
func NewLambdaAPI(session *session.Session, resourceTags *ResourceTags) (*LambdaAPI, error) {
	var api = LambdaAPI{resourceTags: resourceTags}
	// Create Lambda client
	api.client = lambda.New(session)
	return &api, nil
}

func (api *LambdaAPI) doCreateLambda(lambdaName string, nestorID string, roleArn string, task *reporter.Task) (*LambdaInformation, error) {
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
		Runtime:      aws.String(lambda.RuntimeNodejs10X),
		Role:         aws.String(roleArn),
	}
	result, err := api.client.CreateFunction(input)
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
		functionName: *result.FunctionName,
		arn:          *result.FunctionArn,
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
		functionName: *result.Configuration.FunctionName,
		arn:          *result.Configuration.FunctionArn,
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

func (api *LambdaAPI) createLambda(lambdaName string, nestorID string, roleArn string, task *reporter.Task) (*LambdaInformation, error) {
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
			"functionName": lambdaInformation.functionName,
			"arn":          lambdaInformation.arn,
		})

		return lambdaInformation, nil
	}

	t2 := t0.Sub("lambda does not exist - creating it")
	result, err := api.doCreateLambda(lambdaName, nestorID, roleArn, t2)
	if err != nil {
		t2.Fail(err)
	}
	t2.Ok()
	t0.Okr(map[string]string{
		"functionName": result.functionName,
		"arn":          result.arn,
	})
	return result, nil
}
