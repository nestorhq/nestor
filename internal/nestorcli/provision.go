package nestorcli

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/awsapi"
	"github.com/nestorhq/nestor/internal/config"
)

// CliProvision processing for provision CLI command
func CliProvision(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.App.Name
	fmt.Printf("Provision:\n")
	fmt.Printf(" environment: %s\n", environment)
	fmt.Printf(" appName    : %s\n", appName)
	fmt.Printf("config: %v\n", nestorConfig)

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
	for _, table := range nestorConfig.Resources.DynamoDbTable {
		var tableName = appName + "-" + environment + "-" + table.ID
		fmt.Printf("tableName: %s\n", tableName)
		api.CreateMonoTable(tableName)
	}
}
