package terraform

import (
	"blowe-cli/lib/terraform/resource"
	"errors"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
)

type ChangeResource interface {
	GetResourceName() string
	GetResourceId() string
}

func ResolveResource(change *tfjson.ResourceChange) (ChangeResource, error) {
	switch change.Type {
	case "aws_security_group_rule":
		return resource.ResolveAwsSecurityGroupRule(change), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown Resource %s", change.Type))
	}
}
