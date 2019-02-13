package cmd

import (
	"os"
	"fmt"

	"github.com/hqhs/gosupport/internal"
	"github.com/spf13/cobra"
)

var options = internal.Options{}

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
	serveCmd.PersistentFlags().StringVarP(
		&options.Domain, "domain", "d", "localhost", "Domain to use in auth links etc.")
	serveCmd.PersistentFlags().StringVarP(
		&options.Port, "port", "p", "8080", "Port to use with provided domain.")
	serveCmd.PersistentFlags().StringVar(
		&options.EmailServer, "smtp-server", "smtp.gmail.com:587",
		"URL to smtp server with port.")
	serveCmd.PersistentFlags().StringVar(
		&options.EmailAddress, "smtp-address", "admin@gmail.com",
		"Email address used to authenticate on server.")
	serveCmd.PersistentFlags().StringVar(
		&options.EmailPassword, "smtp-password", "s3cr3tpwd",
		"Password for smtp authentication")
	rootCmd.AddCommand(serveCmd)
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
