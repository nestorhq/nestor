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

	api, err := awsapi.NewAwsAPI("sls", "us-west-1", "us-west-2")

	if err != nil {
		panic(err)
	}

	// 1: user pool
	var userPoolName = appName + "-" + environment
	api.CreateUserPool(userPoolName)

	// 2: dynamodb tables
	for _, table := range nestorConfig.Resources.DynamoDbTable {
		var tableName = appName + "-" + environment + "-" + table.Id
		fmt.Printf("tableName: %s\n", tableName)
		api.CreateMonoTable(tableName)
	}
}
