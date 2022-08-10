/*
Copyright © 2022 Bryce Lowe <blowe@patreon.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "terraform subcommands",
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	terraformCmd.PersistentFlags().String("terraformdir", "./", "the terraform base directory")
}
