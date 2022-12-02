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
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type User struct{}

func (u User) Analyze(node *parser.Node, line Line) []error {
	errs := []error{}
	if strings.EqualFold(node.Value, "root") {
		errs = append(errs, fmt.Errorf(`USER directive set to root %s could cause an unexpected behavior. In OpenShift, containers are run using arbitrarily assigned user ID`, PrintLineInfo(line)))
	}
	return errs
}
