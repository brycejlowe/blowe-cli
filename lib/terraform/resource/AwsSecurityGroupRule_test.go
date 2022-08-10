package resource

import (
	"blowe-cli/test"
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAwsSecurityGroupRuleChange_GetResourceId(t *testing.T) {
	securityGroupRuleChanges, err := getSecurityGroupRuleJson()
	if err != nil {
		t.Errorf("Error Unmarshalling Test Data: %v", err)
	}

	securityGroupRuleChange := ResolveAwsSecurityGroupRule(securityGroupRuleChanges.ResourceChanges[0])
	expectedResourceId := "sg-1234_ingress_tcp_3000_3001_192.168.1.5/32"
	if securityGroupRuleChange.GetResourceId() != expectedResourceId {
		t.Errorf(
			"Mismatched Resource ID - Got: '%s' but Expected: '%s'",
			securityGroupRuleChange.GetResourceId(),
			expectedResourceId,
		)
	}
}

func TestAwsSecurityGroupRuleChange_GetResourceName(t *testing.T) {
	securityGroupRuleChanges, err := getSecurityGroupRuleJson()
	if err != nil {
		t.Errorf("Error Unmarshalling Test Data: %v", err)
	}

	securityGroupRuleChange := ResolveAwsSecurityGroupRule(securityGroupRuleChanges.ResourceChanges[0])
	expectedResourceName := "module.security_group_test2.aws_security_group_rule.ingress_cidrs[\"192.168.1.5/32-3000-3001\"]"
	if securityGroupRuleChange.GetResourceName() != expectedResourceName {
		t.Errorf(
			"Mismatched Resource Name - Got: '%s' but Expected: '%s'",
			securityGroupRuleChange.GetResourceName(),
			expectedResourceName,
		)
	}

}

func getSecurityGroupRuleJson() (*tfjson.Plan, error) {
	jsonFile, err := os.Open(filepath.Join(test.GetRoot(), "testdata", "security_group_rule_change.json"))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	jsonBytes, _ := ioutil.ReadAll(jsonFile)

	var securityGroupRuleChanges *tfjson.Plan
	if err := json.Unmarshal(jsonBytes, &securityGroupRuleChanges); err != nil {
		return nil, err
	}

	return securityGroupRuleChanges, nil
}
