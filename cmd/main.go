package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var appVersion = "development"

func Run() {
	app := cli.NewApp()
	app.Name = "vyconfigure"
	app.Version = appVersion
	app.Usage = "Declarative configuration for VyOS."
	app.EnableBashCompletion = true
	app.Authors = []*cli.Author{
		{Name: "Charlie Haley", Email: "charlie-haley@users.noreply.github.com"},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "host",
			Usage:   "The hostname of the VyOS HTTP API.",
			EnvVars: []string{"VYCONFIGURE_HOST"},
		},
		&cli.StringFlag{
			Name:    "api-key",
			Usage:   "API key for the HTTP API.",
			EnvVars: []string{"VYCONFIGURE_API_KEY"},
		},
		&cli.StringFlag{
			Name: "config-dir", Value: ".",
			Usage:   "Directory where config is stored.",
			EnvVars: []string{"VYCONFIGURE_CONFIG_DIR"},
		},
		&cli.BoolFlag{
			Name:    "insecure",
			Usage:   "Whether to skip verifying the SSL certificate.",
			EnvVars: []string{"VYCONFIGURE_INSECURE"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Enable Debug mode.",
			EnvVars: []string{"VYCONFIGURE_DEBUG"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name: "version", Aliases: []string{"v"}, Usage: "prints the current version.",
			Action: version,
		},
		{
			Name: "apply", Aliases: []string{"a", "push"}, Usage: "applies the current configuration.",
			Action: apply,
		},
		{
			Name: "sync", Aliases: []string{"s", "pull"}, Usage: "syncs configuration to the current directory through the HTTP API.",
			Action: sync,
		},
		{
			Name: "plan", Aliases: []string{"d", "diff"}, Usage: "shows the diff between the current directory and VyOS instance",
			Action: plan,
		},
		{
			Name: "show", Aliases: []string{"sh"}, Usage: "shows the local instance",
			Action: show,
		},
	}
	app.Action = version

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
