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
	fmt.Printf("Provision:\n")
	fmt.Printf(" environment: %s\n", environment)
	fmt.Printf(" appName    : %s\n", appName)
	fmt.Printf("config: %v\n", nestorConfig)

	var _ = reporter.NewReporterM(reporter.NewMessage("command: provision").WithArg("environment", environment)).Start()

	// TODO: hard coded
	resourceTags := awsapi.NewResourceTag("1", environment, appName)

	// TODO: hard coded
	api, err := awsapi.NewAwsAPI("sls", resourceTags, "us-west-1", "us-west-2")

	if err != nil {
		panic(err)
	}

	// 1: user pool
	var userPoolName = appName + "-" + environment
	upr, errup := api.CreateUserPool(userPoolName)
	if errup != nil {
		panic(errup)
	}
	fmt.Printf("user pool: %v\n", upr)

	// 2: dynamodb tables
	var tableName = appName + "-" + environment + "-main"
	fmt.Printf("tableName: %s\n", tableName)
	api.CreateMonoTable(tableName)
}
