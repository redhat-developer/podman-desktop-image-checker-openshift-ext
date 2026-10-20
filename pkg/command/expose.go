/**********************************************************************
 * Copyright (C) 2022 Red Hat, Inc.
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
	if strings.HasPrefix(str, "map[") && strings.HasSuffix(str, "]") {
		str = str[4 : len(str)-1]
	}
	index := strings.IndexByte(str, '/')
	if index >= 0 {
		str = str[0:index]
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
