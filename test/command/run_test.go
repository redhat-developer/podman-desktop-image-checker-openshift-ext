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
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/command"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

/*
chown -R node:node /app
chown --recursive=node:node
chown -R $ZOOKEEPER_USER:$HADOOP_GROUP $ZOOKEEPER_LOG_DIR
chown -R 1000:1000 /app
chown 1001 /deployments/run-java.sh
chown -h 501:20 './AirRun Updates'
*/
func TestCorrectParsingOfChownCommandWithUserAndRootGroup(t *testing.T) {
	verifyParsingCommand(t, "chown -R node:root /app", 0)
}

func TestCorrectParsingOfChownCommandWithUserAndRootGroupAndLongFlag(t *testing.T) {
	verifyParsingCommand(t, "chown --recursive=node:root /app", 0)
}

func TestCorrectParsingOfChownCommandWithUserAndRootGroupAsNumber(t *testing.T) {
	verifyParsingCommand(t, "chown -R 1000:0 /app", 0)
}

func TestCorrectParsingOfChownCommandWithUserAndRootGroupAsNumberAndLongFlag(t *testing.T) {
	verifyParsingCommand(t, "chown --recursive=1000:0 /app", 0)
}

func TestFailIfChownCommandWithUserAndNonRootGroup(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chown -R node:node /app", 1)
	if !strings.Contains(suggestions[0].Error(), "In OpenShift the group ID must always be set to the root group (0)") {
		t.Errorf("Expected to be wrong group ID error but it was %s", suggestions[0].Error())
	}
}

func TestFailIfChownCommandWithUserAndNonRootGroupAndLongFlag(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chown --recursive=node:node /app", 1)
	if !strings.Contains(suggestions[0].Error(), "In OpenShift the group ID must always be set to the root group (0)") {
		t.Errorf("Expected to be wrong group ID error but it was %s", suggestions[0].Error())
	}
}

func TestFailIfChownCommandWithUserAndNonRootGroupAsNumber(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chown -R 1000:1000 /app", 1)
	if !strings.Contains(suggestions[0].Error(), "In OpenShift the group ID must always be set to the root group (0)") {
		t.Errorf("Expected to be wrong group ID error but it was %s", suggestions[0].Error())
	}
}

func TestFailIfChownCommandWithUserAndNonRootGroupAsNumberAndLongFlag(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chown --recursive=1000:1000 /app", 1)
	if !strings.Contains(suggestions[0].Error(), "In OpenShift the group ID must always be set to the root group (0)") {
		t.Errorf("Expected to be wrong group ID error but it was %s", suggestions[0].Error())
	}
}

func TestCorrectParsingOfChownCommandWithOnlyUserSet(t *testing.T) {
	verifyParsingCommand(t, "chown -R node /app", 0)
}

func TestCorrectChmodCommandWithExecuteGroupPermission(t *testing.T) {
	verifyParsingCommand(t, "chmod 070 /app", 0)
}

func TestFailChmodCommandWithNonGroupPermission(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chmod 000 /app", 1)
	if !strings.Contains(suggestions[0].Error(), "permission set on") {
		t.Errorf("Expected to be wrong group permissions error but it was %s", suggestions[0].Error())
	}
}

func TestFailChmodCommandWithUserSetButNotGroup(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chmod 700 /app", 1)
	if !strings.Contains(suggestions[0].Error(), "permission set on") {
		t.Errorf("Expected to be wrong group permissions error but it was %s", suggestions[0].Error())
	}
}

func TestFailChmodCommandWithInvalidPermissionCode(t *testing.T) {
	suggestions := verifyParsingCommand(t, "chmod 70 /app", 1)
	if !strings.Contains(suggestions[0].Error(), "unable to fetch args of chmod command") {
		t.Errorf("Expected to be unable to fetch args of chmod command but it was %s", suggestions[0].Error())
	}
}

func verifyParsingCommand(t *testing.T, cmd string, numberExpectedErrors int) []error {
	run := command.Run{}
	suggestions := run.Analyze(&parser.Node{
		Value: cmd,
	},
		utils.Source{
			Name: "test",
			Type: utils.Image,
		},
		command.Line{
			Start: 1,
			End:   1,
		})
	if len(suggestions) != numberExpectedErrors {
		t.Errorf("Expected %d suggestions but they were %d", numberExpectedErrors, len(suggestions))
	}
	return suggestions
}
