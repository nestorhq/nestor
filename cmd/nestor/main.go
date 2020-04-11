package main

import (
	"fmt"
	"os"

	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/nestorcli"
	"github.com/nestorhq/nestor/internal/reporter"
	"github.com/segmentio/cli"
)

func main() {
	reporter.Experiment()
	os.Exit(0)

	nestorConfig, err := config.ReadConfig("nestor.yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("app: %s\n", nestorConfig.App.Name)

	type cliEnvConfig struct {
		Environment string `flag:"-e,--environment" help:"environment to use" default:"dev"`
	}

	cliProvision := cli.Command(func(config cliEnvConfig) {
		nestorcli.CliProvision(config.Environment, nestorConfig)
	})

	cli.Exec(cli.CommandSet{
		"provision": cliProvision,
	})
}
