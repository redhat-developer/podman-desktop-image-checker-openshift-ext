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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type Expose struct {
}

type exposeResultKeyType struct{}

var exposeResultKey exposeResultKeyType

func (e Expose) Analyze(ctx context.Context, node *parser.Node, source utils.Source, line Line) context.Context {
	str := node.Value
	index := strings.IndexByte(node.Value, '/')
	if index >= 0 {
		str = node.Value[0:index]
	}
	var results []Result
	port, err := strconv.Atoi(str)
	if err != nil {
		results = append(results, Result{
			Name:        "Wrong port value",
			Status:      StatusFailed,
			Severity:    SeverityCritical,
			Description: err.Error(),
		})
	}
	if port < 1024 {
		results = append(results, Result{
			Name:        "Privileged port exposed",
			Status:      StatusFailed,
			Severity:    SeverityHigh,
			Description: fmt.Sprintf(`port %d exposed %s could be wrong. TCP/IP port numbers below 1024 are privileged port numbers`, port, GenerateErrorLocation(source, line)),
		})
	}
	return context.WithValue(ctx, exposeResultKey, results)
}

func (e Expose) PostProcess(ctx context.Context) []Result {
	result := ctx.Value(exposeResultKey)
	if result == nil {
		return nil
	}
	return result.([]Result)
}
