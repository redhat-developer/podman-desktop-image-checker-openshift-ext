package command

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type User struct{}

func (u User) Analyze(node *parser.Node, line Line) []error {
	errs := []error{}
	if strings.EqualFold(node.Value, "root") {
		errs = append(errs, fmt.Errorf(`USER directive set to root %s could cause an unexpected behavior. In OpenShift, containers are run using arbitrarily assigned user ID`, PrintLineInfo(line)))
	}
	return errs
}
