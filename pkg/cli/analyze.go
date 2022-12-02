package cli

import (
	"fmt"

	analyzer "github.com/redhat-developer/docker-openshift-analizer/pkg/command"
	"github.com/spf13/cobra"
)

func NewCmdAnalyze() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:     "analyze",
		Short:   "Analyze the Dockerfile and discover potential issues when deploying it on OpenShift",
		Long:    "Analyze the Dockerfile and discover potential issues when deploying it on OpenShift",
		Args:    cobra.MaximumNArgs(1),
		Run:     doAnalyze,
		Example: `  doa analyze /your/local/project/path`,
	}
	return analyzeCmd
}

func doAnalyze(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		PrintNoArgsWarningMessage(cmd.Name())
		return
	}

	PrintPrettifyOutput(analyzer.AnalyzePath(args[0]))
}

func PrintNoArgsWarningMessage(command string) {
	fmt.Printf(`
No arg received. Did you forget to add the dockerfile or project path to analyze?

Expected:
  doa %s /your/local/project/path [flags]

To find out more, run 'doa %s --help'
`, command, command)
}

func PrintPrettifyOutput(errs []error) {
	for i, sug := range errs {
		fmt.Printf("%d - %s\n\n", i+1, sug)
	}
}
