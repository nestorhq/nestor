package actions

import (
	"errors"

	"github.com/nestorhq/nestor/internal/reporter"
)

// description of static lambda information
type lambdaConfigInformation struct {
	functioName          string
	environmentVariables map[string]string
}

func (actions *Actions) findLambdaInformationFromID(ID string) (*lambdaConfigInformation, error) {
	var nestorConfig = actions.nestorConfig
	var appName = actions.nestorConfig.Application.Name
	var environment = actions.environment
	for _, lambdaFunction := range nestorConfig.Resources.LambdaFunction {
		if lambdaFunction.ID == ID {
			// TODO: refactor to have a single method defining the lambda name
			var lambdaFunctionName = appName + "-" + environment + "-" + lambdaFunction.ID
			var environmentVariables = map[string]string{}
			for _, envVar := range lambdaFunction.Environment {
				environmentVariables[envVar.Name] = envVar.Value
			}
			var result = &lambdaConfigInformation{
				functioName:          lambdaFunctionName,
				environmentVariables: environmentVariables,
			}
			return result, nil
		}
	}
	return nil, errors.New("No lambda found with id:" + ID)
}

// DoDeployment perform deploymenst
func (actions *Actions) DoDeployment(task *reporter.Task) error {
	var appName = actions.nestorConfig.Application.Name
	var environment = actions.environment
	var api = actions.api
	var nestorConfig = actions.nestorConfig
	// var nestorResources = actions.nestorResources

	var t = task.SubM(reporter.NewMessage("DoDeployment").
		WithArg("environment", environment).
		WithArg("appName", appName))

	for _, lambdaDeployment := range nestorConfig.Deployments.LambdaFunction {
		t0 := t.SectionM(reporter.NewMessage("Deploy lambda").
			WithArg("id", lambdaDeployment.ID).
			WithArg("file", lambdaDeployment.File).
			WithArg("handler", lambdaDeployment.Handler))

		lambdaIformation, err := actions.findLambdaInformationFromID(lambdaDeployment.ID)
		if err != nil {
			t0.Fail(err)
			return err
		}
		zipFileName := lambdaDeployment.File
		handler := lambdaDeployment.Handler
		err = api.UpdateLambdaCodeFromZip(lambdaIformation.functioName, zipFileName, handler,
			lambdaIformation.environmentVariables, t0)
		if err != nil {
			t0.Fail(err)
			return err
		}
	}

	return nil
}
