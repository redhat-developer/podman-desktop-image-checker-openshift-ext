/*******************************************************************************
 * Copyright (c) 2022 Red Hat, Inc.
 * Distributed under license by Red Hat, Inc. All rights reserved.
 * This program is made available under the terms of the
 * Eclipse Public License v2.0 which accompanies this distribution,
 * and is available at http://www.eclipse.org/legal/epl-v20.html
 *
 * Contributors:
 * Red Hat, Inc.
 ******************************************************************************/
package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type Line struct {
	Start int
	End   int
}
type Command interface {
	Analyze(*parser.Node, utils.Source, Line) []error
}

var commandHandlers = map[string]Command{
	utils.EXPOSE_INSTRUCTION: Expose{},
	utils.FROM_INSTRUCTION:   From{},
	utils.RUN_INSTRUCTION:    Run{},
	utils.USER_INSTRUCTION:   User{},
}

func AnalyzePath(path string) []error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return []error{fmt.Errorf("unable to analyze %s - error %s", path, err)}
	}

	if fileInfo.IsDir() {
		path = filepath.Join(path, "Dockerfile")
	}

	file, err := os.Open(path)
	if err != nil {
		return []error{fmt.Errorf("unable to open %s - error %s", path, err)}
	}
	defer file.Close()

	return AnalyzeFile(file)
}

func AnalyzeFile(file *os.File) []error {
	res, err := parser.Parse(file)
	if err != nil {
		return []error{
			fmt.Errorf("unable to analyze the Dockerfile. Error when parsing %s : %s", file.Name(), err.Error()),
		}
	}

	return AnalyzeNodeFromSource(res.AST, utils.Source{
		Name: "",
		Type: utils.Image,
	})
}

func AnalyzeNodeFromSource(node *parser.Node, source utils.Source) []error {
	suggestions := []error{}
	commands := []string{}
	for _, child := range node.Children {
		commands = append(commands, child.Original) // TODO to be used if we need to check previous rows to make sugestions
		line := Line{
			Start: child.StartLine,
			End:   child.EndLine,
		}
		handler := commandHandlers[strings.ToUpper(child.Value+" ")]
		if handler != nil {
			for n := child.Next; n != nil; n = n.Next {
				if n.Value == "" {
					suggestions = append(suggestions, fmt.Errorf("%s %s has an empty value", child.Value, GenerateErrorLocation(source, line)))
				} else {
					suggestions = append(suggestions, handler.Analyze(n, source, line)...)
				}
			}
		}
	}
	return suggestions
}

func IsCommand(text string, command string) bool {
	return strings.Contains(text, command)
}

func GenerateErrorLocation(source utils.Source, line Line) string {
	if source.Type == utils.Parent {
		return fmt.Sprintf("in parent image %s", source.Name)
	}
	if line.Start == line.End {
		return fmt.Sprintf("at line %d", line.Start)
	}
	return fmt.Sprintf("at line %d-%d", line.Start, line.End)
}
