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
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Run struct{}

func (r Run) Analyze(node *parser.Node) []error {
	errs := []error{}

	// let's split the run command by &&. E.g chmod 070 /app && chmod 070 /app/routes && chmod 070 /app/bin
	splittedCommands := strings.Split(node.Value, "&&")
	for _, command := range splittedCommands {
		if r.isChmodCommand(command) {
			err := r.analyzeChmodCommand(command)
			if err != nil {
				errs = append(errs, err)
			}
		} else if r.isChownCommand(command) {
			err := r.analyzeChownCommand(command)
			if err != nil {
				errs = append(errs, err)
			}
		}

	}

	return errs
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
func (r Run) analyzeChownCommand(s string) error {
	re := regexp.MustCompile(`(\$*\w+)*:(\$*\w+)`)

	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil // errors.New("unable to find any group set by the chown command")
	}
	group := match[len(match)-1]
	if strings.ToLower(group) != "root" && group != "0" {
		return fmt.Errorf(`warning owner set on %s .
		In OpenShift the group ID must always be set to the root group (0)`, s)
	}
	return nil
}

func (r Run) isChmodCommand(s string) bool {
	return IsCommand(s, "chmod")
}

func (r Run) analyzeChmodCommand(s string) error {
	re := regexp.MustCompile(`chmod\s+(\d+)\s+(.*)`)
	match := re.FindStringSubmatch(s)
	if len(match) != 3 {
		return errors.New("unable to fetch args of chmod command")
	}
	permission := match[1]
	if len(permission) != 3 {
		return errors.New("unable to fetch args of chmod command")
	}
	groupPermission := permission[1:2]
	if groupPermission != "7" {
		return fmt.Errorf(`warning permission set on %s .
		In OpenShift, the directories and files that the processes running in the image need to access 
		should have their group ownership set to the root group. They also need to be read/writable by that group as 
		recommended by the OpenShift Container Platform-specific guidelines`, s)
	}
	return nil
}
