package utils

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"testing"
)

func TestExtractSourceDestination(t *testing.T) {
	ipv4Cidrs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, ipv4Cidr := range ipv4Cidrs {
		result := ExtractSourceDestination(types.SecurityGroupRule{
			CidrIpv4: &ipv4Cidr,
		})

		if result != ipv4Cidr {
			t.Error("Error Matching IPv4 of Security Group Rule")
		}
	}

	ipv6Cidrs := []string{
		"2001:4860:4860::8888/125",
		"2001:4860:4860::8890/124",
		"2001:4860:4860::88a0/123",
	}

	for _, ipv6Cidr := range ipv6Cidrs {
		result := ExtractSourceDestination(types.SecurityGroupRule{
			CidrIpv4: &ipv6Cidr,
		})

		if result != ipv6Cidr {
			t.Error("Error Matching IPv6 of Security Group Rule")
		}
	}

}
