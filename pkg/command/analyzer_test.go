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
