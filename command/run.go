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
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Run struct{}

func (r Run) Analyze(node *parser.Node, line Line) []error {
	errs := []error{}

	// let's split the run command by &&. E.g chmod 070 /app && chmod 070 /app/routes && chmod 070 /app/bin
	splittedCommands := strings.Split(node.Value, "&&")
	for _, command := range splittedCommands {
		if r.isChmodCommand(command) {
			err := r.analyzeChmodCommand(command, line)
			if err != nil {
				errs = append(errs, err)
			}
		} else if r.isChownCommand(command) {
			err := r.analyzeChownCommand(command, line)
			if err != nil {
				errs = append(errs, err)
			}
		} else if r.isSudoOrSuCommand(command) {
			err := r.analyzeSudoAndSuCommand(command, line)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func (r Run) isSudoOrSuCommand(s string) bool {
	return IsCommand(s, "sudo") || IsCommand(s, "su")
}

func (r Run) analyzeSudoAndSuCommand(s string, line Line) error {
	re := regexp.MustCompile(`(\s+|^)(sudo|su)\s+`)

	match := re.FindStringSubmatch(s)
	if len(match) > 0 {
		return fmt.Errorf(`sudo/su command used in '%s' %s could cause an unexpected behavior. 
		In OpenShift, containers are run using arbitrarily assigned user ID and elevating privileges could lead 
		to runtime errors`, s, PrintLineInfo(line))
	}
	return nil
}

func (r Run) isChownCommand(s string) bool {
	return IsCommand(s, "chown")
}

/* to be tested on
chown -R node:node /app
chown --recursive=node:node
chown +x test
RUN chown -R $ZOOKEEPER_USER:$HADOOP_GROUP $ZOOKEEPER_LOG_DIR
chown -R 1000:1000 /app
chown 1001 /deployments/run-java.sh
chown -h 501:20 './AirRun Updates'
*/
func (r Run) analyzeChownCommand(s string, line Line) error {
	re := regexp.MustCompile(`(\$*\w+)*:(\$*\w+)`)

	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil // errors.New("unable to find any group set by the chown command")
	}
	group := match[len(match)-1]
	if strings.ToLower(group) != "root" && group != "0" {
		return fmt.Errorf(`owner set on %s %s could cause an unexpected behavior. 
		In OpenShift the group ID must always be set to the root group (0)`, s, PrintLineInfo(line))
	}
	return nil
}

func (r Run) isChmodCommand(s string) bool {
	return IsCommand(s, "chmod")
}

func (r Run) analyzeChmodCommand(s string, line Line) error {
	re := regexp.MustCompile(`chmod\s+(\d+)\s+(.*)`)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil
	}
	if len(match) != 3 {
		return fmt.Errorf("unable to fetch args of chmod command %s. Is it correct?", PrintLineInfo(line))
	}
	permission := match[1]
	if len(permission) != 3 {
		return fmt.Errorf("unable to fetch args of chmod command %s. Is it correct?", PrintLineInfo(line))
	}
	groupPermission := permission[1:2]
	if groupPermission != "7" {
		proposal := fmt.Sprintf("Is it an executable file? Try updating permissions to %s7%s", permission[0:1], permission[2:3])
		if groupPermission != "6" {
			proposal += fmt.Sprintf(" otherwise set it to %s6%s", permission[0:1], permission[2:3])
		}
		return fmt.Errorf("permission set on %s %s could cause an unexpected behavior. %s\n"+
			"Explanation - in Openshift, directories and files need to be read/writable by the root group and "+
			"files that must be executed should have group execute permissions", s, PrintLineInfo(line), proposal)
	}

	return nil
}
