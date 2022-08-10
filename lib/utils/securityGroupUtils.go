package utils

import "github.com/aws/aws-sdk-go-v2/service/ec2/types"

func ExtractSourceDestination(input types.SecurityGroupRule) string {
	sourceDestination := ""
	if input.CidrIpv6 != nil {
		sourceDestination = *input.CidrIpv6
	}

	if input.CidrIpv4 != nil {
		sourceDestination = *input.CidrIpv4
	}

	if input.ReferencedGroupInfo != nil {
		sourceDestination = *input.ReferencedGroupInfo.GroupId
	}

	return sourceDestination
}
