package utils

import (
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/redhat-developer/docker-openshift-analyzer/pkg/utils"
	"strings"
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

func Line2Node(line string, root *parser.Node) error {
	result, err := parser.Parse(strings.NewReader(line))
	if err != nil {
		return err
	}
	for _, node := range result.AST.Children {
		root.AddChild(node, -1, -1)
	}
	return nil
}

func ExtractCmd(str string) string {
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
