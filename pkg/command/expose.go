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
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type Expose struct {
}

var exposeUuid = uuid.New()

func (e Expose) UUID() uuid.UUID {
	return exposeUuid
}

func (e Expose) Analyze(ctx AnalyzeContext, node *parser.Node, source utils.Source, line Line) {
	str := node.Value
	commandContext := ctx.CommandContext(exposeUuid)
	index := strings.IndexByte(node.Value, '/')
	if index >= 0 {
		str = node.Value[0:index]
	}
	port, err := strconv.Atoi(str)
	if err != nil {
		commandContext.Results = append(commandContext.Results, Result{
			Name:        "Wrong port value",
			Status:      StatusFailed,
			Severity:    SeverityCritical,
			Description: err.Error(),
		})
	}
	if port < 1024 {
		commandContext.Results = append(commandContext.Results, Result{
			Name:        "Privileged port exposed",
			Status:      StatusFailed,
			Severity:    SeverityHigh,
			Description: fmt.Sprintf(`port %d exposed %s could be wrong. TCP/IP port numbers below 1024 are privileged port numbers`, port, GenerateErrorLocation(source, line)),
		})
	}
}

func (e Expose) PostProcess(context AnalyzeContext) {
}
