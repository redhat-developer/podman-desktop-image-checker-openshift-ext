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
)

type Expose struct {
}

func (e Expose) Analyze(node *parser.Node) []error {
	errs := []error{}
	for n := node.Next; n != nil; n = n.Next {
		port, err := strconv.Atoi(n.Value)
		if err != nil {
			errs = append(errs, err)
		}
		if port < 1024 {
			errs = append(errs, fmt.Errorf(`dockerfile exposes port %d. TCP/IP port numbers below 1024 are privileged port numbers 
			that enable only the root user to bind to these ports`, port))
		}
	}
	return errs
}
