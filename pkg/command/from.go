/*******************************************************************************
 * Copyright (c) 2023 Red Hat, Inc.
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
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type From struct {
}

type fromResultKeyType struct{}

var fromResultKey fromResultKeyType

const SCRATCH_IMAGE_NAME = "scratch"

func (f From) Analyze(ctx context.Context, node *parser.Node, source utils.Source, line Line) context.Context {
	if node.Value == SCRATCH_IMAGE_NAME {
		return ctx
	}
	decompiledNode, err := decompiler.Decompile(node.Value)
	if err != nil {
		// unable to decompile base image
		ctx = context.WithValue(ctx, fromResultKey, []Result{
			Result{
				Name:        "Analyze error",
				Status:      StatusFailed,
				Severity:    SeverityLow,
				Description: fmt.Sprintf("unable to analyze the base image %s", node.Value),
			},
		})
	}
	_, ctx = AnalyzeNodeFromSource(ctx, decompiledNode, utils.Source{
		Name: node.Value,
		Type: utils.Parent,
	})
	return ctx
}

func (f From) PostProcess(ctx context.Context) []Result {
	result := ctx.Value(fromResultKey)
	if result == nil {
		return nil
	}
	return result.([]Result)
}
