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
package main

import (
	"flag"
	"os"

	"github.com/redhat-developer/docker-openshift-analyzer/pkg/cli"
)

func main() {
	doaCmd := cli.DockerOpenShiftAnalyzerCommands()
	flag.Usage = func() {
		_ = doaCmd.Help()
	}
	// parse the flags but hack around to avoid exiting with error code 2 on help
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	args := os.Args[1:]
	if err := flag.CommandLine.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if err := doaCmd.Execute(); err != nil {
		cli.RedirectErrorStringToStdErrAndExit(err.Error())
	}
}
