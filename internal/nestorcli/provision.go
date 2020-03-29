package nestorcli

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/awsApi"
	"github.com/nestorhq/nestor/internal/config"
)

func CliProvision(environment string, nestorConfig *config.Config) {
	fmt.Printf("CliProvision: Environment is %s\n", environment)
	fmt.Printf("config: %v\n", nestorConfig)

	_, err := awsApi.NewAwsApi("sls")

	if err != nil {
		panic(err)
	}
}
