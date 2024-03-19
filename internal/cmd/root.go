package cmd

import (
	"runtime"

	"github.com/ChrisWiegman/kana/internal/console"
	"github.com/ChrisWiegman/kana/internal/site"

	"github.com/spf13/cobra"
)

var (
	flagName                    string
	flagVerbose, flagJSONOutput bool
	commandsRequiringSite       []string
)

func Execute() {
	kanaSite := new(site.Site)
	consoleOutput := new(console.Console)

	// Setup the cobra command
	cmd := &cobra.Command{
		Use:   "kana",
		Short: "Kana is a simple WordPress development tool designed for plugin and theme developers.",
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			consoleOutput.Debug = flagVerbose
			consoleOutput.JSON = flagJSONOutput

			if cmd.Use == "wp" {
				err := parseWPNameFlag(args, cmd)
				if err != nil {
					consoleOutput.Error(err)
				}
			}

			err := kanaSite.New(cmd, commandsRequiringSite, startFlags, flagVerbose, consoleOutput, Version)
			if err != nil {
				consoleOutput.Error(err)
			}
		},
	}

	// Hide the default completion command
	cmd.CompletionOptions.HiddenDefaultCmd = true

	// Add the "name" flag to allow for sites not connected to the local directory
	cmd.PersistentFlags().StringVarP(&flagName, "name", "n", "", "Specify a name for the site, used to override using the current folder.")
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Display debugging information along with detailed command output")
	cmd.PersistentFlags().BoolVar(&flagJSONOutput, "output-json", false, "Display all output in JSON format for further processing")

	err := cmd.PersistentFlags().MarkHidden("output-json")
	if err != nil {
		consoleOutput.Error(err)
	}

	// Register the subcommands
	cmd.AddCommand(
		changelog(consoleOutput),
		config(consoleOutput, kanaSite),
		db(consoleOutput, kanaSite),
		destroy(consoleOutput, kanaSite),
		export(consoleOutput, kanaSite),
		flush(consoleOutput, kanaSite),
		list(consoleOutput, kanaSite),
		open(consoleOutput, kanaSite),
		start(consoleOutput, kanaSite),
		stop(consoleOutput, kanaSite),
		version(consoleOutput),
		wp(consoleOutput, kanaSite),
		xdebug(consoleOutput, kanaSite),
	)

	if runtime.GOOS == "darwin" {
		cmd.AddCommand(trust(consoleOutput, kanaSite))
	}

	// Execute anything we need to
	if err := cmd.Execute(); err != nil {
		consoleOutput.Error(err)
	}
}
