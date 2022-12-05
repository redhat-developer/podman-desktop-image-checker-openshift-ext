package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	analyzer "github.com/redhat-developer/docker-openshift-analyzer/pkg/command"
	"github.com/spf13/cobra"
)

func NewCmdAnalyze() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:     "analyze",
		Short:   "Analyze the Dockerfile and discover potential issues when deploying it on OpenShift",
		Long:    "Analyze the Dockerfile and discover potential issues when deploying it on OpenShift. It accepts the project root path or the Dockerfile path.",
		Args:    cobra.MaximumNArgs(1),
		Run:     doAnalyze,
		Example: `  doa analyze /your/local/project/path[/Dockerfile_name]`,
	}
	analyzeCmd.PersistentFlags().String(
		"o", "", "Specify output format, supported format: json",
	)
	return analyzeCmd
}

func doAnalyze(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		PrintNoArgsWarningMessage(cmd.Name())
		return
	}

	out := cmd.Flag("o")
	if out.Value.String() != "" && !strings.EqualFold(out.Value.String(), "json") {
		RedirectErrorStringToStdErrAndExit(fmt.Sprintf("unknown value '%s' for flag %s, type --help for a list of all flags\n", out.Value.String(), out.Name))
	} else if strings.EqualFold(out.Value.String(), "json") {
		PrintPrettifyJsonOutput(analyzer.AnalyzePath(args[0]))
		return
	}

	PrintPrettifyOutput(analyzer.AnalyzePath(args[0]))
}

func PrintNoArgsWarningMessage(command string) {
	fmt.Printf(`
No arg received. Did you forget to add the dockerfile or project path to analyze?

Expected:
  doa %s /your/local/project/path[/Dockerfile_name] [flags]

To find out more, run 'doa %s --help'
`, command, command)
}

func PrintPrettifyJsonOutput(errs []error) {
	rowErrs := make([]string, len(errs))
	for i, err := range errs {
		rowErrs[i] = err.Error()
	}
	var bytes []byte
	var err error
	if bytes, err = json.MarshalIndent(rowErrs, "", "    "); err != nil {
		fmt.Println("error while converting output to json. Please try again without the output (--o) flag")
	}
	fmt.Println(string(bytes))
}

func PrintPrettifyOutput(errs []error) {
	for i, sug := range errs {
		fmt.Printf("%d - %s\n\n", i+1, sug)
	}
}
