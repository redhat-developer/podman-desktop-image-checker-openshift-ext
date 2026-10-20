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
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type User struct{}

type userResultKeyType struct{}
type userProcessedKeyType struct{}

var userResultKey userResultKeyType
var userProcessedKey userProcessedKeyType

func (u User) Analyze(ctx context.Context, node *parser.Node, source utils.Source, line Line) context.Context {
	var results []Result
	if strings.EqualFold(node.Value, "root") {
		results = append(results, Result{
			Name:        "User set to root",
			Status:      StatusFailed,
			Severity:    SeverityMedium,
			Description: fmt.Sprintf(`USER directive set to root %s could cause an unexpected behavior. In OpenShift, containers are run using arbitrarily assigned user ID`, GenerateErrorLocation(source, line)),
		})
	}
	ctx = context.WithValue(ctx, userResultKey, results)
	return context.WithValue(ctx, userProcessedKey, true)
}

func (u User) PostProcess(ctx context.Context) []Result {
	processed := ctx.Value(userProcessedKey)
	res := ctx.Value(userResultKey)
	var results []Result
	if res != nil {
		results = res.([]Result)
	}
	if processed == nil {
		results = append(results, Result{
			Name:        "User set to root",
			Status:      StatusFailed,
			Severity:    SeverityMedium,
			Description: fmt.Sprintf("USER directive implicitely set to root could cause an unexpected behavior. In OpenShift, containers are run using arbitrarily assigned user ID"),
		})
	}
	return results

}
