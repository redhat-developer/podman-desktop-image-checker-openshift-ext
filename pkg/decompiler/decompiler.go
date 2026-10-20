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
 package decompiler

import (
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/pkg/errors"
	docker "github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler/docker"
	podman "github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler/podman"
	registry "github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler/registry"
)

type Provider interface {
	Decompile(imageName string) (*parser.Node, error)
}

func Decompile(imageName string) (*parser.Node, error) {
	providers := []Provider{
		podman.PodmanProvider{},
		docker.DockerProvider{},
		registry.RegistryProvider{},
	}
	for _, provider := range providers {
		node, err := provider.Decompile(imageName)
		if err != nil {
			return nil, err
		}
		if node != nil {
			return node, nil
		}
	}
	return nil, errors.Errorf("Can't resolve image %s", imageName)
}
