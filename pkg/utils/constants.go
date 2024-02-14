/**********************************************************************
 * Copyright (C) 2024 Red Hat, Inc.
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
 package utils

const NOP = "#(nop) "

const RUN_PREFIX = "/bin/sh -c "

const FROM_INSTRUCTION = "FROM "
const RUN_INSTRUCTION = "RUN "
const CMD_INSTRUCTION = "CMD "
const LABEL_INSTRUCTION = "LABEL "
const MAINTAINER_INSTRUCTION = "MAINAINER "
const EXPOSE_INSTRUCTION = "EXPOSE "
const ENV_INSTRUCTION = "ENV "
const ADD_INSTRUCTION = "ADD "
const COPY_INSTRUCTION = "COPY "
const ENTRYPOINT_INSTRUCTION = "ENTRYPOINT "
const VOLUME_INSTRUCTION = "VOLUME "
const USER_INSTRUCTION = "USER "
const WORKDIR_INSTRUCTION = "WORKDIR "
const ARG_INSTRUCTION = "ARG "
const ONBUILD_INSTRUCTION = "ONBUILD "
const STOPSIGNAL_INSTRUCTION = "STOPSIGNAL "
const HEALTHCHECK_INSTRUCTION = "HEALTHCHECK "
const SHELL_INSTRUCTION = "SHELL "

type SourceType string

const (
	Image  SourceType = "IMAGE"
	Parent SourceType = "PARENT_IMAGE"
)

type Source struct {
	Name string
	Type SourceType
}
