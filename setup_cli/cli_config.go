package setup_cli

import (
	"os"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/codegangsta/cli"
	"github.com/pivotal-cf-experimental/diego-edge-cli/app_runner"
	"github.com/pivotal-cf-experimental/diego-edge-cli/config"
	"github.com/pivotal-cf-experimental/diego-edge-cli/config/config_helpers"
	"github.com/pivotal-cf-experimental/diego-edge-cli/config/persister"

	app_runner_command_factory "github.com/pivotal-cf-experimental/diego-edge-cli/app_runner/command_factory"
	config_command_factory "github.com/pivotal-cf-experimental/diego-edge-cli/config/command_factory"
)

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Diego"
	app.Usage = "Command line interface for diego."

	config := config.New(persister.NewFilePersister(config_helpers.ConfigFileLocation(userHome())))
	config.Load()

	receptorClient := receptor.NewClient(config.Api(), "", "")
	appRunner := app_runner.NewDiegoAppRunner(receptorClient)

	app.Commands = []cli.Command{
		app_runner_command_factory.NewAppRunnerCommandFactory(appRunner, os.Stdout).MakeStartDiegoAppCommand(),
		app_runner_command_factory.NewAppRunnerCommandFactory(appRunner, os.Stdout).MakeScaleDiegoAppCommand(),
		config_command_factory.NewConfigCommandFactory(config, os.Stdout).MakeSetTargetCommand(),
	}
	return app
}

func userHome() string {
	return os.Getenv("HOME")
}
