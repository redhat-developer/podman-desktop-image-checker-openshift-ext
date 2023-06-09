package decompiler

import (
	"sort"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	decompilerutils "github.com/redhat-developer/docker-openshift-analyzer/pkg/decompiler/utils"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

type OrderedHistory []v1.History

func (o OrderedHistory) Len() int {
	return len(o)
}

func (o OrderedHistory) Less(i, j int) bool {
	return o[i].Created.Before(o[j].Created.Time)
}

func (o OrderedHistory) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type RegistryProvider struct{}

func (p RegistryProvider) Decompile(imageName string) (*parser.Node, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, nil
	}

	configFile, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}

	root := &parser.Node{}

	history := configFile.History
	sort.Sort(OrderedHistory(history))
	for _, hist := range history {
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

	if configFile.Config.User != "" {
		err := decompilerutils.Line2Node(utils.USER_INSTRUCTION+configFile.Config.User, root)
		if err != nil {
			return nil, err
		}
	}

	return root, nil
}
