package nestorcli

import (
	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/reporter"
)

// CliDeploy processing for deploy CLI command
func CliDeploy(environment string, nestorConfig *config.Config) {
	var appName = nestorConfig.Application.Name

	var t = reporter.NewReporterM(reporter.NewMessage("command: deploy").
		WithArg("environment", environment).
		WithArg("appName", appName)).
		Start()
	t.Log("Deploy command -- TODO")
	t.Ok()
}
