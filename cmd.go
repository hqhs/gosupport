package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gosupport",
	Short: "Management commands for support bot dashboard",
	Long:  "From top to bottom re-written version of usefull tool for in-ompany helpdesking.",
	Run: func(cmd *cobra.Command, args []string) {
		// run webserver
	},
}

// ExecuteCmd executes given command and exit if recieved en error
func ExecuteCmd() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
