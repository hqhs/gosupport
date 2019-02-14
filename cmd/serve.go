package cmd

import (
	"fmt"
	"os"

	"github.com/hqhs/gosupport/internal/app"
	"github.com/spf13/cobra"
	"github.com/hqhs/gosupport/pkg/templator"
	kitlog "github.com/go-kit/kit/log"
)

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
	serveCmd.PersistentFlags().StringVar(
		&options.DatabaseURL, "database url", "", "Url to connect to get database access")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Simple tech-support dashboard for in-company usage",
	Long: `Yet bare-bones dashboard for managing chats with anyone, who needs help,
 with telegrambot interface`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// TODO allow user set logger level (better to do it globally)
		l := kitlog.NewLogfmtLogger(os.Stdout)
		l.Log("root", options.Root, "address", options.Domain + ":" + options.Port)
		t, err := templator.NewTemplator(options.Root)
		if err != nil {
			// NOTE we dont panic here to allow init process to finish and find all
			// errors, since they're independent. But http router serving won't start.
			l.Log("panic", err)
		}
		l.Log("templates", fmt.Sprintf("%+v", t.GetTemplates()))
		m := app.NewMockMailer(t, l)
		db := app.NewMockDatabase()
		if err != nil {
			l.Log("status", "exiting", "message", "Fix 'panic' errors above to serve http requests")
			os.Exit(1)
		}
		s := app.InitServer(l, t, m, db, options)
		// TODO start polling messages from bots
		s.ServeRouter()
	},
}
