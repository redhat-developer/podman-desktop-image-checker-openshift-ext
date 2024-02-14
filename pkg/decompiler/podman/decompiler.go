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
	"context"
	"github.com/containers/common/pkg/config"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/images"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"sort"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	decompilerutils "github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler/utils"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type OrderedHistory []v1.History

func (o OrderedHistory) Len() int {
	return len(o)
}

func (o OrderedHistory) Less(i, j int) bool {
	return o[i].Created.Before(*o[j].Created)
}

func (o OrderedHistory) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type PodmanProvider struct{}

func getPodmanConnection() (string, string) {
	conf, err := config.NewConfig("")
	if err != nil {
		return "", ""
	}
	uri := ""
	identity := ""
	if conf.Engine.ActiveService != "" {
		uri = conf.Engine.ServiceDestinations[conf.Engine.ActiveService].URI
		identity = conf.Engine.ServiceDestinations[conf.Engine.ActiveService].Identity
	}
	return uri, identity
}

func (p PodmanProvider) Decompile(imageName string) (*parser.Node, error) {
	uri, identity := getPodmanConnection()
	if uri != "" {
		ctx, err := bindings.NewConnectionWithIdentity(context.Background(), uri, identity, false)
		if err != nil {
			return nil, nil
		}
		image, err := images.GetImage(ctx, imageName, nil)
		if err != nil {
			return nil, nil
		}

		root := &parser.Node{}
		sort.Sort(OrderedHistory(image.History))
		for _, hist := range image.History {
			if hist.Comment != "" && strings.HasPrefix(strings.ToUpper(hist.Comment), utils.FROM_INSTRUCTION) &&
				!hist.EmptyLayer {
				err := decompilerutils.Line2Node(hist.Comment, root)
				if err != nil {
					return nil, err
				}
			}
			if hist.CreatedBy != "" {
				cmd := decompilerutils.ExtractCmd(hist.CreatedBy)
				if cmd != "" {
					err := decompilerutils.Line2Node(cmd, root)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		return root, nil
	}
	return nil, nil
}
