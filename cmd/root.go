package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
)

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "admiral_sharding_managersr",
		Short: "Admiral Sharding Manager is a load distributor",
		Long:  "Admiral Sharding Manager manages load distribution among admiral operators",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Admiral Sharding Manager")
		},
	}
	rootCmd.SetArgs(args)
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	return rootCmd
}
