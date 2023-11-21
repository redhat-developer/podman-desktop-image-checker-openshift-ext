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

import "testing"

func TestCheckNginx(t *testing.T) {
	for _, tag := range []string{"1.25.0", "1.25.1", "1.25.2", "1.25.3"} {
		t.Run(tag, func(t *testing.T) {
			AnalyzeImage("docker.io/nginx:" + tag)
		})
	}
}

func TestFromScratch(t *testing.T) {
	errors := AnalyzePath("resources/Containerfile.fromscratch")
	if len(errors) != 1 {
		t.Error("Image with FROM scratch returns errors")
	}
}
func TestFromNginxWithUser(t *testing.T) {
	errors := AnalyzePath("resources/Containerfile.fromnginxwithuser")
	if len(errors) != 1 {
		t.Error("Image with FROM nginx with USER returns errors")
	}
}
