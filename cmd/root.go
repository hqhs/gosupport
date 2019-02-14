package cmd

import (
	"os"
	"fmt"

	"github.com/hqhs/gosupport/internal/app"
	"github.com/spf13/cobra"
)

var options = app.Options{}

var rootCmd = &cobra.Command{
	Use:   "support",
	Short: "Simple tech-support dashboard for in-company usage",
	Long: `Bare-bones dashboard for managing chats with anyone, who needs help,
 with telegrambot interface`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	options.DbOptions = app.DbOptions{}
}

// Execute executes cli commands
func Execute() {
	if root, err := os.Getwd(); err != nil {
		fmt.Printf("Couldn't stat current directory, %v\n", err)
		os.Exit(1)
	} else {
		options.Root = root
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
