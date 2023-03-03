package decompiler

import (
	"context"
	"github.com/containers/common/pkg/config"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
)

var CONTAINERFILE_INSTRUCTIONS = []string{
	utils.RUN_INSTRUCTION,
	utils.CMD_INSTRUCTION,
	utils.LABEL_INSTRUCTION,
	utils.MAINTAINER_INSTRUCTION,
	utils.EXPOSE_INSTRUCTION,
	utils.ENV_INSTRUCTION,
	utils.ADD_INSTRUCTION,
	utils.COPY_INSTRUCTION,
	utils.ENTRYPOINT_INSTRUCTION,
	utils.VOLUME_INSTRUCTION,
	utils.USER_INSTRUCTION,
	utils.WORKDIR_INSTRUCTION,
	utils.ARG_INSTRUCTION,
	utils.ONBUILD_INSTRUCTION,
	utils.STOPSIGNAL_INSTRUCTION,
	utils.HEALTHCHECK_INSTRUCTION,
	utils.SHELL_INSTRUCTION,
}

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

		for _, hist := range image.History {
			if hist.Comment != "" && strings.HasPrefix(strings.ToUpper(hist.Comment), utils.FROM_INSTRUCTION) &&
				!hist.EmptyLayer {
				err := line2Node(hist.Comment, root)
				if err != nil {
					return nil, err
				}
			}
			if hist.CreatedBy != "" {
				cmd := extractCmd(hist.CreatedBy)
				if cmd != "" {
					err := line2Node(cmd, root)
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

func line2Node(line string, root *parser.Node) error {
	result, err := parser.Parse(strings.NewReader(line))
	if err != nil {
		return err
	}
	for _, node := range result.AST.Children {
		root.AddChild(node, -1, -1)
	}
	return nil
}

func extractCmd(str string) string {
	index := strings.Index(str, utils.NOP)
	if index > 0 {
		return str[index+len(utils.NOP):]
	}
	index = strings.Index(str, utils.RUN_PREFIX)
	if index >= 0 {
		return utils.RUN_INSTRUCTION + str[index+len(utils.RUN_PREFIX):]
	}
	if isContainerFileInstruction(str) {
		return str
	}
	return ""
}

func isContainerFileInstruction(str string) bool {
	for _, prefix := range CONTAINERFILE_INSTRUCTIONS {
		if strings.HasPrefix(strings.ToUpper(str), prefix) {
			return true
		}
	}
	return false
}
