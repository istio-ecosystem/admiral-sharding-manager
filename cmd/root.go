/*
Copyright © 2024 Intuit Inc.
*/
package cmd

import (
	"log"

	"github.com/istio-ecosystem/admiral-sharding-manager/pkg/model"

	"os"

	"github.com/spf13/cobra"
)

// maintains sharding manger program arguments
var smParams = model.ShardingManagerParams{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "admiral-sharding-manager",
	Short: "Sharding manager distributes load among admiral operators",
	Long:  "Sharding manager distributes load among admiral operators",
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("admiral sharding manager has been initialized")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// manage root command flags
func init() {}
