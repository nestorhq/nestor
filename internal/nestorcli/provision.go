package nestorcli

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/awsApi"
	"github.com/nestorhq/nestor/internal/config"
)

func CliProvision(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.App.Name
	fmt.Printf("Provision:\n")
	fmt.Printf(" environment: %s\n", environment)
	fmt.Printf(" appName    : %s\n", appName)
	fmt.Printf("config: %v\n", nestorConfig)

	api, err := awsApi.NewAwsApi("sls", "us-west-1")

	if err != nil {
		panic(err)
	}

	// 1: dynamodb tables
	for _, table := range nestorConfig.Resources.DynamoDbTable {
		var tableName = appName + "-" + environment + "-" + table.Id
		fmt.Printf("tableName: %s\n", tableName)
		api.CreateMonoTable(tableName)
	}
}
