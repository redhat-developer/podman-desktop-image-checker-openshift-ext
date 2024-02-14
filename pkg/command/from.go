/**********************************************************************
 * Copyright (C) 2023 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 ***********************************************************************/
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
		return context.WithValue(ctx, fromResultKey, []Result{
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
