package main

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/nestorcli"
	"github.com/segmentio/cli"
)

func main() {
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
