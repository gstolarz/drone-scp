package main

import (
	"github.com/urfave/cli/v2"

	"github.com/gstolarz/drone-scp/plugin"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "address",
			Usage:       "address",
			EnvVars:     []string{"PLUGIN_ADDRESS"},
			Destination: &settings.Address,
		},
		&cli.StringFlag{
			Name:        "username",
			Usage:       "username",
			EnvVars:     []string{"PLUGIN_USERNAME"},
			Destination: &settings.Username,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "password",
			EnvVars:     []string{"PLUGIN_PASSWORD"},
			Destination: &settings.Password,
		},
		&cli.StringFlag{
			Name:        "key",
			Usage:       "key",
			EnvVars:     []string{"PLUGIN_KEY"},
			Destination: &settings.Key,
		},
		&cli.StringFlag{
			Name:        "source",
			Usage:       "source",
			EnvVars:     []string{"PLUGIN_SOURCE"},
			Destination: &settings.Source,
		},
		&cli.StringFlag{
			Name:        "target",
			Usage:       "target",
			EnvVars:     []string{"PLUGIN_TARGET"},
			Destination: &settings.Target,
		},
		&cli.BoolFlag{
			Name:        "templating",
			Usage:       "templating",
			EnvVars:     []string{"PLUGIN_TEMPLATING"},
			Destination: &settings.Templating,
		},
	}
}
