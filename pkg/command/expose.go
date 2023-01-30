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
	"strconv"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type Expose struct {
}

func (e Expose) Analyze(node *parser.Node, source utils.Source, line Line) []error {
	errs := []error{}
	port, err := strconv.Atoi(node.Value)
	if err != nil {
		errs = append(errs, err)
	}
	if port < 1024 {
		errs = append(errs, fmt.Errorf(`port %d exposed %s could be wrong. TCP/IP port numbers below 1024 are privileged port numbers`, port, GenerateErrorLocation(source, line)))
	}
	return errs
}
