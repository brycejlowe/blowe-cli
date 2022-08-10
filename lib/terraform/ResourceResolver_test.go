package terraform

import (
	tfjson "github.com/hashicorp/terraform-json"
	"testing"
)

func TestResolveResource(t *testing.T) {
	// test valid resource
	validResource, err := ResolveResource(&tfjson.ResourceChange{
		Type: "aws_security_group_rule",
	})

	if validResource == nil {
		t.Error("Valid Resource: Unexpected nil")
	}

	if err != nil {
		t.Error("Valid Resource: Expected nil Error")
	}

	// test invalid resource
	invalidResource, err := ResolveResource(&tfjson.ResourceChange{
		Type: "i'm-invalid",
	})

	if invalidResource != nil {
		t.Error("Invalid Resource: Needs nil Resource")
	}

	if err == nil {
		t.Error("Invalid Resource: Needs Error")
	}
}
