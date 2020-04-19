package main

import (
	"fmt"

	"github.com/nestorhq/nestor/internal/config"
	"github.com/nestorhq/nestor/internal/nestorcli"
	"github.com/segmentio/cli"
)

func main() {
	// reporter.Experiment()
	// os.Exit(0)
	fmt.Println("@@ nestor")
	nestorConfig, err := config.ReadConfig("nestor.yml")
	if err != nil {
		panic(err)
	}

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
