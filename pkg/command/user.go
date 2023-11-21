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
	"github.com/google/uuid"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type User struct{}

var userUuid = uuid.New()

const PROCESSED_KEY = "processed"

func (u User) UUID() uuid.UUID {
	return userUuid
}

func (u User) Analyze(ctx AnalyzeContext, node *parser.Node, source utils.Source, line Line) {
	commandContext := ctx.CommandContext(userUuid)
	if strings.EqualFold(node.Value, "root") {
		commandContext.Results = append(commandContext.Results, Result{
			Name:        "User set to root",
			Status:      StatusFailed,
			Severity:    SeverityMedium,
			Description: fmt.Sprintf(`USER directive set to root %s could cause an unexpected behavior. In OpenShift, containers are run using arbitrarily assigned user ID`, GenerateErrorLocation(source, line)),
		})
	} else {
		commandContext.Results = nil
	}
	commandContext.Infos[PROCESSED_KEY] = true
}

func (u User) PostProcess(ctx AnalyzeContext) {
	commandContext := ctx.CommandContext(userUuid)
	if _, processed := commandContext.Infos[PROCESSED_KEY]; !processed {
		commandContext.Results = append(commandContext.Results, Result{
			Name:        "User set to root",
			Status:      StatusFailed,
			Severity:    SeverityMedium,
			Description: fmt.Sprintf("USER directive implicitely set to root could cause an unexpected behavior. In OpenShift, containers are run using arbitrarily assigned user ID"),
		})
	}

}
