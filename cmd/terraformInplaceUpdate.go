/*
Copyright © 2022 Bryce Lowe <blowe@patreon.com>

*/
package cmd

import (
	"blowe-cli/lib/terraform"
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// inplaceUpdateCmd represents the inplaceUpdate command
var inplaceUpdateCmd = &cobra.Command{
	Use:   "inplaceUpdate",
	Short: "Update security group rules in-place.",
	Long: `This command relies on a valid terraform plan and essentially does state operations in a loop, rather than
having Terraform do potentially destructive adds and deletes.  It will take a binary plan that's been pumped through
terraform show json-izer (i.e. terraform show -json ../foo.tfplan) and do the terraform state rm and terraform import
for you.

Note: I've left out concurrence here as the commands will essentially be serialized behind writing to a single 
json file.  I also discovered that there's no clear Terraform SDK for doing these operations, even in the go client.
The recommendation seems to be shelling out to the terraform executable which is what I ended up doing.
`,
	Run: func(cmd *cobra.Command, args []string) {
		terraformDirectory, _ := cmd.Flags().GetString("terraformdir")
		planFilePath, _ := cmd.Flags().GetString("jsonplan")

		tfCommand := terraform.NewCmd(terraformDirectory)
		version, err := tfCommand.GetVersion()

		if err != nil {
			log.Fatalf("Error Fetching Terraform Version: %s\n", err.(*exec.ExitError).Stderr)
		}

		log.Printf("Terraform Version: %s\n", version)

		planFile, err := os.Open(planFilePath)
		if err != nil {
			log.Fatalln("Error Opening File")
		}

		fileBytes, err := ioutil.ReadAll(planFile)
		if err != nil {
			log.Fatalln("Error Reading File Contents")
		}

		var tfPlanResourceChanges tfjson.Plan
		if err := json.Unmarshal(fileBytes, &tfPlanResourceChanges); err != nil {
			log.Fatalf("Error Parsing Terraform Plan Json: %v\n", err)
		}

		deleteQueue := make([]terraform.ChangeResource, 0)
		createQueue := make([]terraform.ChangeResource, 0)
		for _, resourceChange := range tfPlanResourceChanges.ResourceChanges {
			for _, changeAction := range resourceChange.Change.Actions {
				if changeAction == tfjson.ActionCreate {
					createResource, err := terraform.ResolveResource(resourceChange)
					if err != nil {
						log.Fatalf("Error Fetching Create Resource %s", resourceChange.Address)
					}
					createQueue = append(createQueue, createResource)
					log.Printf("Queuing Create for %s", resourceChange.Address)
				}

				if changeAction == tfjson.ActionDelete {
					deleteResource, err := terraform.ResolveResource(resourceChange)
					if err != nil {
						log.Fatalf("Error Fetching Create Resource %s", resourceChange.Address)
					}
					deleteQueue = append(deleteQueue, deleteResource)
					log.Printf("Queuing Delete for %s", resourceChange.Address)
				}
			}
		}

		log.Printf("Delete: %d, Import: %d", len(deleteQueue), len(createQueue))

		log.Println("Processing Delete Queue")
		deleteProcessedCount := 0
		deleteErrorCount := 0
		for _, val := range deleteQueue {
			deleteProcessedCount++
			log.Printf("Delete: '%s' [%d of %d]", val.GetResourceName(), deleteProcessedCount, len(deleteQueue))
			if err := tfCommand.DoStateRm(val.GetResourceName()); err != nil {
				deleteErrorCount++
				log.Printf("Error:\n%s", string(err.(*exec.ExitError).Stderr))
			}
		}
		log.Printf("Finished Processing Delete Queue with %d Failures", deleteErrorCount)

		log.Println("Processing Import Queue")
		createProcessedCount := 0
		createErrorCount := 0
		for _, val := range createQueue {
			createProcessedCount++
			log.Printf("Import: '%s' using id '%s' [%d of %d]", val.GetResourceName(), val.GetResourceId(), createProcessedCount, len(createQueue))
			if err := tfCommand.DoImport(val.GetResourceName(), val.GetResourceId()); err != nil {
				createErrorCount++
				log.Printf("Error:\n%s", string(err.(*exec.ExitError).Stderr))
			}
		}
		log.Printf("Finished Processing Import Queue with %d Failures", createErrorCount)

		log.Printf("Delete Error(s): %d, Import Error(s): %d", deleteErrorCount, createErrorCount)
		if deleteErrorCount+createErrorCount > 0 {
			os.Exit(1)
		}

		log.Println("Complete")
	},
}

func init() {
	terraformCmd.AddCommand(inplaceUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inplaceUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inplaceUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	inplaceUpdateCmd.Flags().String("jsonplan", "", "path to the json plan")

	_ = inplaceUpdateCmd.MarkFlagRequired("jsonplan")
}
