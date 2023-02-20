package decompiler

import (
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/pkg/errors"
	registry "github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler/registry"
)

type Provider interface {
	Decompile(imageName string) (*parser.Node, error)
}

func Decompile(imageName string) (*parser.Node, error) {
	providers := []Provider{
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
	return nil, errors.Errorf("Can't resolved umage %s", imageName)
}
