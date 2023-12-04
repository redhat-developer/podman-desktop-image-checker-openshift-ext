package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	doaLong = `
The OpenShift Image Checker is a CLI tool for finding and highlighting potential issues a Containerfile could have on an OpenShift cluster.
Find out more at https://github.com/redhat-developer/podman-desktop-image-checker-ext	`

	doaExample = `
  # Analyze the Containerfile of a project:
    doa analyze /your/local/project/path[/Containerfile_name]
	`

	rootHelpMessage = "To see a full list of commands, run 'doa --help'"

	rootDefaultHelp = fmt.Sprintf("%s\n\nExamples:\n%s\n\n%s", doaLong, doaExample, rootHelpMessage)
)

func DockerOpenShiftAnalyzerCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "doa",
		Short:   "doa",
		Long:    doaLong,
		RunE:    ShowHelp,
		Example: doaExample,
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Create a custom help function that will exit when we enter an invalid command, for example:
	// doa foobar --help
	// which will exit with an error message: "unknown command 'foobar', type --help for a list of all commands"
	helpCmd := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(command *cobra.Command, args []string) {
		// Simple way of checking to see if the command has a parent (if it doesn't, it does not exist)
		if !command.HasParent() && len(args) > 0 {
			RedirectErrorStringToStdErrAndExit(fmt.Sprintf("unknown command '%s', type --help for a list of all commands\n", args[0]))
		}
		helpCmd(command, args)
	})

	rootCmdList := append([]*cobra.Command{},
		NewCmdAnalyze(),
	)

	rootCmd.AddCommand(rootCmdList...)

	return rootCmd
}

func RedirectErrorStringToStdErrAndExit(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// ShowHelp will show the help correctly (and whether or not the command is invalid...)
// Taken from: https://github.com/redhat-developer/odo/blob/f55a4f0a7af4cd5f7c4e56dd70a66d38be0643cf/pkg/odo/cli/cli.go#L272
func ShowHelp(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		// We will show a custom help when typing JUST `doa`, directing the user to use `doa --help` for a full help.
		// Thus we will set cmd.SilenceUsage and cmd.SilenceErrors both to true so we do not output the usage or error out.
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		// Print out the default "help" usage
		fmt.Println(rootDefaultHelp)
		return nil
	}

	//revive:disable:error-strings This is a top-level error message displayed as is to the end user
	return errors.New("invalid command - see available commands/subcommands above")
	//revive:enable:error-strings
}
