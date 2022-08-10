/*
Copyright © 2022 Bryce Lowe <blowe@patreon.com>

*/
package cmd

import (
	"blowe-cli/lib/terraform"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

// terraformPlanJsonCmd represents the terraformPlanJson command
var terraformPlanJsonCmd = &cobra.Command{
	Use:   "planJson",
	Short: "generate a json plan",
	Long:  `Use Terraform's built-in functionality to generate a JSON based plan.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating Terraform Json Plan")
		terraformDirectory, _ := cmd.Flags().GetString("terraformdir")
		planPath, _ := cmd.Flags().GetString("plan")
		jsonPlanPath, _ := cmd.Flags().GetString("jsonplan")

		tfCommand := terraform.NewCmd(terraformDirectory)

		log.Println("Running Terraform Plan")
		if err := tfCommand.GeneratePlan(planPath); err != nil {
			log.Fatalf("Error Running Plan:\n%s", err.(*exec.ExitError).Stderr)
		}

		log.Println("Converting Terraform Plan to json")
		if err := tfCommand.ConvertPlanToJson(planPath, jsonPlanPath); err != nil {
			log.Fatalf("Error Converting Plan:\n%s", err.(*exec.ExitError).Stderr)
		}

		log.Println("Complete")
	},
}

func init() {
	terraformCmd.AddCommand(terraformPlanJsonCmd)

	terraformPlanJsonCmd.Flags().String("plan", "./plan.tfplan", "temporary file path")
	terraformPlanJsonCmd.Flags().String("jsonplan", "./plan.tfplan.json", "json plan output")
}
