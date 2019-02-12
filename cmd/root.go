package cmd

import (
	"os"
	"fmt"

	"github.com/hqhs/gosupport/app"

	"github.com/spf13/cobra"
)

var options = app.Options{}

var rootCmd = &cobra.Command{
	Use:   "support",
	Short: "Simple tech-support dashboard for in-company usage",
	Long: `Bare-bones dashboard for managing chats with anyone, who needs help,
 with telegrambot interface`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world!")
		// s := app.InitServer()
		// Dashboard = InitSite(Env)
		// if err := Dashboard.InitDatabase(); err != nil {
		// 	fmt.Println("Something went wrong during database initialization: ", err)
		// 	os.Exit(1)
		// }
		// if err := Dashboard.InitBots(); err != nil {
		// 	fmt.Println("Something went wrong during bots initialization: ", err)
		// 	os.Exit(1)
		// }
		// if err := Dashboard.InitializeRouter(); err != nil {
		// 	fmt.Println("Something went wrong during router initialization: ", err)
		// 	os.Exit(1)
		// }
		// Dashboard.RunBots()
		// Dashboard.RunSite()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&options.Domain, "domain", "d", "localhost", "Domain to use in auth links etc.")
	rootCmd.PersistentFlags().StringVarP(
		&options.Port, "port", "p", "8080", "Port to use with provided domain.")
	rootCmd.PersistentFlags().StringVar(
		&options.EmailServer, "smtp-server", "smtp.gmail.com:587",
		"URL to smtp server with port.")
	rootCmd.PersistentFlags().StringVar(
		&options.EmailAddress, "smtp-address", "admin@gmail.com",
		"Email address used to authenticate on server.")
	rootCmd.PersistentFlags().StringVar(
		&options.EmailPassword, "smtp-password", "s3cr3tpwd",
		"Password for smtp authentication")
}

// Execute executes cli commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
