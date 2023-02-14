package decompiler

import (
	"sort"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
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

func Decompile(imageName string) (*parser.Node, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}

	configFile, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}

	root := &parser.Node{}

	history := configFile.History
	sort.Sort(OrderedHistory(history))
	processed := false
	for _, hist := range history {
		if hist.CreatedBy != "" {
			cmd := extractCmd(hist.CreatedBy)
			if cmd != "" {
				err := line2Node(cmd, root)
				if err != nil {
					return nil, err
				}
				processed = true
			}
		}
		if !processed && hist.Comment != "" && strings.HasPrefix(strings.ToUpper(hist.Comment), utils.FROM_INSTRUCTION) {
			err := line2Node(hist.Comment, root)
			if err != nil {
				return nil, err
			}
		}
	}

	if configFile.Config.User != "" {
		err := line2Node(utils.USER_INSTRUCTION+configFile.Config.User, root)
		if err != nil {
			return nil, err
		}
	}

	return root, nil
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
