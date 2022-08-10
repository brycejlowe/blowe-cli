/*
Copyright © 2022 Bryce Lowe <blowe@patreon.com>

*/
package cmd

import (
	"blowe-cli/lib"
	"blowe-cli/lib/utils"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	"log"
	"sort"
	"strings"
)

// listSecurityGroupsCmd represents the listSecurityGroups command
var listSecurityGroupsCmd = &cobra.Command{
	Use:   "listSecurityGroups",
	Short: "List security groups matching a source/destination",
	Long: `List the security group names of the groups that match a particular source/destination.

Currently the source and/or destination is matched using a contains cause that's all I need right now.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Listing Security Groups")

		sourceDestinationValue, _ := cmd.Flags().GetString("value")
		if len(sourceDestinationValue) > 0 {
			log.Printf("Search Criteria: %s", sourceDestinationValue)
		}

		awsConfig := lib.GetConfig()
		client := ec2.NewFromConfig(awsConfig)

		log.Println("Fetching Security Group Rules")
		rulesInput := &ec2.DescribeSecurityGroupRulesInput{}
		rulesPaginator := ec2.NewDescribeSecurityGroupRulesPaginator(client, rulesInput)
		foundGroupRules := make(map[string][]types.SecurityGroupRule)
		for rulesPaginator.HasMorePages() {
			output, err := rulesPaginator.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("Error Fetching Rules: %v", err)
			}

			for _, r := range output.SecurityGroupRules {
				// no search criteria passed
				if len(sourceDestinationValue) == 0 {
					foundGroupRules[*r.GroupId] = append(foundGroupRules[*r.GroupId], r)
					continue
				}

				if strings.Contains(utils.ExtractSourceDestination(r), sourceDestinationValue) {
					foundGroupRules[*r.GroupId] = append(foundGroupRules[*r.GroupId], r)
				}
			}
		}

		log.Printf("Found %d Groups with Matching Source/Destination Value", len(foundGroupRules))

		// batch the security group ids so we can call the api more efficiently
		chunks := 10
		log.Printf("Batching Security Groups in Groups of %d", chunks)
		var groups [][]string
		var group []string
		for k, _ := range foundGroupRules {
			group = append(group, k)
			if len(group) == chunks {
				groups = append(groups, group)
				group = nil
			}
		}

		// catch anything that didn't end up in a group
		if len(group) > 0 {
			groups = append(groups, group)
			group = nil
		}

		log.Println("Fetching Security Group Information from AWS")
		foundGroups := make(map[string]types.SecurityGroup)
		for _, groupId := range groups {
			// fetch the security group header information
			groupsInput := &ec2.DescribeSecurityGroupsInput{
				Filters: []types.Filter{
					{
						Name:   aws.String("group-id"),
						Values: groupId,
					},
				},
			}

			groupsPaginator := ec2.NewDescribeSecurityGroupsPaginator(client, groupsInput)
			for groupsPaginator.HasMorePages() {
				result, err := groupsPaginator.NextPage(context.TODO())
				if err != nil {
					log.Fatalf("Error Fetching Groups: %v", err)
				}

				for _, g := range result.SecurityGroups {
					foundGroups[*g.GroupId] = g
				}
			}
		}

		log.Println("** Start Group Information **")
		var groupNames []string
		for _, v := range foundGroups {
			groupNames = append(groupNames, *v.GroupName)
		}

		sort.Strings(groupNames)
		for _, v := range groupNames {
			fmt.Println(v)
		}

		log.Println("** End Group Information **")
	},
}

func init() {
	awsCmd.AddCommand(listSecurityGroupsCmd)

	listSecurityGroupsCmd.Flags().String("value", "", "rule value to search for")
}
