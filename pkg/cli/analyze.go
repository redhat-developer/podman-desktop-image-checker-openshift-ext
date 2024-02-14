/**********************************************************************
 * Copyright (C) 2024 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 ***********************************************************************/

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
		Short:   "Analyze the Containerfile and discover potential issues when deploying it on OpenShift",
		Long:    "Analyze the Containerfile and discover potential issues when deploying it on OpenShift. It accepts the project root path or the Containerfile path.",
		Args:    cobra.MaximumNArgs(0),
		Run:     doAnalyze,
		Example: `  doa analyze -f /your/local/project/path[/Containerfile_name]`,
	}
	analyzeCmd.PersistentFlags().StringP(
		"file", "f", "", "Container file to analyze",
	)
	analyzeCmd.PersistentFlags().StringP(
		"image", "i", "", "Image name to analyze",
	)
	analyzeCmd.PersistentFlags().StringP(
		"output", "o", "", "Specify output format, supported format: json",
	)
	return analyzeCmd
}

func doAnalyze(cmd *cobra.Command, args []string) {
	containerfile := cmd.Flag("file")
	image := cmd.Flag("image")
	if containerfile.Value.String() == "" && image.Value.String() == "" {
		PrintNoArgsWarningMessage(cmd.Name())
		return
	}

	outputFunc := PrintPrettifyOutput
	out := cmd.Flag("output")
	if out.Value.String() != "" && !strings.EqualFold(out.Value.String(), "json") {
		RedirectErrorStringToStdErrAndExit(fmt.Sprintf("unknown value '%s' for flag %s, type --help for a list of all flags\n", out.Value.String(), out.Name))
	} else if strings.EqualFold(out.Value.String(), "json") {
		outputFunc = PrintPrettifyJsonOutput
	}

	if containerfile.Value.String() != "" {
		outputFunc(analyzer.AnalyzePath(containerfile.Value.String()))
	} else if image.Value.String() != "" {
		outputFunc(analyzer.AnalyzeImage(image.Value.String()))
	}
}

func PrintNoArgsWarningMessage(command string) {
	fmt.Printf(`
No arg received. Did you forget to add the Containerfile or project path to analyze?

Expected:
  doa %s /your/local/project/path[/Containerfile_name] [flags]

To find out more, run 'doa %s --help'
`, command, command)
}

func PrintPrettifyJsonOutput(results []analyzer.Result) {
	var bytes []byte
	var err error
	if bytes, err = json.MarshalIndent(results, "", "    "); err != nil {
		fmt.Println("error while converting output to json. Please try again without the output (--o) flag")
	}
	fmt.Println(string(bytes))
}

func PrintPrettifyOutput(results []analyzer.Result) {
	for i, sug := range results {
		fmt.Printf("%d - %s (%s): %s\n\n", i+1, sug.Name, sug.Severity, sug.Description)
	}
}
