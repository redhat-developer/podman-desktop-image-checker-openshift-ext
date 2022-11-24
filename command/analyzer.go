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
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Command interface {
	Analyze(node *parser.Node) []error
}

var commandHandlers = map[string]Command{
	"expose": Expose{},
	"run":    Run{},
}

func Analyze(file *os.File) []error {
	res, err := parser.Parse(file)
	if err != nil {
		return []error{
			fmt.Errorf("unable to analyze the Dockerfile. Error when parsing %s : %s", file.Name(), err.Error()),
		}
	}

	suggestions := []error{}
	commands := []string{}
	for _, child := range res.AST.Children {
		commands = append(commands, child.Original) // TODO to be used if we need to check previous rows to make sugestions
		handler := commandHandlers[strings.ToLower(child.Value)]
		if handler != nil {
			for n := child.Next; n != nil; n = n.Next {
				if n.Value == "" {
					suggestions = append(suggestions, fmt.Errorf("%s has an empty value", child.Value))
				} else {
					suggestions = append(suggestions, handler.Analyze(n)...)
				}
			}
		}
	}
	return suggestions
}

func IsCommand(text string, command string) bool {
	return strings.Contains(text, command)
}
