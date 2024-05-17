package main

import (
	"github.com/istio-ecosystem/admiral-sharding-manager/cmd"
	"os"
)

func main() {
	rootCmd := cmd.GetRootCmd(os.Args[1:])

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(-1)
	}
}
