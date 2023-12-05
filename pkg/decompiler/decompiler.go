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
