package resource

import (
	"encoding/json"
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"strings"
)

type AutoGenerated struct {
	Before struct {
		CidrBlocks            []string      `json:"cidr_blocks"`
		Description           string        `json:"description"`
		FromPort              int           `json:"from_port"`
		ID                    string        `json:"id"`
		Ipv6CidrBlocks        []string      `json:"ipv6_cidr_blocks"`
		PrefixListIds         []interface{} `json:"prefix_list_ids"`
		Protocol              string        `json:"protocol"`
		SecurityGroupID       string        `json:"security_group_id"`
		Self                  bool          `json:"self"`
		SourceSecurityGroupID interface{}   `json:"source_security_group_id"`
		ToPort                int           `json:"to_port"`
		Type                  string        `json:"type"`
	} `json:"before"`
	After struct {
		CidrBlocks            []string    `json:"cidr_blocks"`
		Description           interface{} `json:"description"`
		FromPort              int         `json:"from_port"`
		Ipv6CidrBlocks        []string    `json:"ipv6_cidr_blocks"`
		PrefixListIds         interface{} `json:"prefix_list_ids"`
		Protocol              string      `json:"protocol"`
		SecurityGroupID       string      `json:"security_group_id"`
		Self                  bool        `json:"self"`
		SourceSecurityGroupID interface{} `json:"source_security_group_id"`
		ToPort                int         `json:"to_port"`
		Type                  string      `json:"type"`
	} `json:"after"`
	AfterUnknown struct {
		CidrBlocks            []bool `json:"cidr_blocks"`
		ID                    bool   `json:"id"`
		SourceSecurityGroupID bool   `json:"source_security_group_id"`
	} `json:"after_unknown"`
	BeforeSensitive struct {
		CidrBlocks     []bool        `json:"cidr_blocks"`
		Ipv6CidrBlocks []interface{} `json:"ipv6_cidr_blocks"`
		PrefixListIds  []interface{} `json:"prefix_list_ids"`
	} `json:"before_sensitive"`
	AfterSensitive struct {
		CidrBlocks []bool `json:"cidr_blocks"`
	} `json:"after_sensitive"`
	ReplacePaths [][]string `json:"replace_paths"`
}

func ResolveAwsSecurityGroupRule(change *tfjson.ResourceChange) *AwsSecurityGroupRuleChange {
	// TODO: this is stupid, I'm just marshal/unmarshaling to cast to the right struct
	resourceChangeRaw, _ := json.Marshal(change.Change)
	var resourceChange *AutoGenerated
	_ = json.Unmarshal(resourceChangeRaw, &resourceChange)

	return &AwsSecurityGroupRuleChange{
		resourceChange:      change,
		securityGroupChange: resourceChange,
	}
}

type AwsSecurityGroupRuleChange struct {
	resourceChange      *tfjson.ResourceChange
	securityGroupChange *AutoGenerated
}

func (c *AwsSecurityGroupRuleChange) GetResourceName() string {
	return c.resourceChange.Address
}

func (c *AwsSecurityGroupRuleChange) GetResourceId() string {
	groupChange := c.securityGroupChange.After

	source := make([]string, 0)

	if groupChange.Self == true {
		source = append(source, "self")
	}

	if groupChange.CidrBlocks != nil {
		source = append(source, groupChange.CidrBlocks...)
	}

	if groupChange.Ipv6CidrBlocks != nil {
		source = append(source, groupChange.Ipv6CidrBlocks...)
	}

	if groupChange.SourceSecurityGroupID != nil {
		source = append(source, groupChange.SourceSecurityGroupID.(string))
	}

	return fmt.Sprintf(
		"%s_%s_%s_%d_%d_%s",
		groupChange.SecurityGroupID,
		groupChange.Type,
		groupChange.Protocol,
		groupChange.FromPort,
		groupChange.ToPort,
		strings.Join(source, "_"),
	)
}
