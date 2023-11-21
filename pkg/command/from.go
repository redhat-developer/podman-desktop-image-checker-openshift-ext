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
	"fmt"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type From struct {
}

const SCRATCH_IMAGE_NAME = "scratch"

func (f From) Analyze(node *parser.Node, source utils.Source, line Line) []Result {
	if node.Value == SCRATCH_IMAGE_NAME {
		return nil
	}
	results := []Result{}
	decompiledNode, err := decompiler.Decompile(node.Value)
	if err != nil {
		// unable to decompile base image
		results = append(results, Result{
			Name:        "Analyze error",
			Status:      StatusFailed,
			Severity:    SeverityLow,
			Description: fmt.Sprintf("unable to analyze the base image %s", node.Value),
		})
		return results
	}
	errsFromBaseImage := AnalyzeNodeFromSource(decompiledNode, utils.Source{
		Name: node.Value,
		Type: utils.Parent,
	})
	if len(errsFromBaseImage) > 0 {
		results = append(results, errsFromBaseImage...)
	}
	return results
}
