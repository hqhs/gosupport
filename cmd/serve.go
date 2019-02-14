package cmd

import (
	"fmt"
	"os"
	"syscall"
	"os/signal"
	"path/filepath"

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
		&options.StaticFiles, "static", "web/static", "Path to directory with static files")
	serveCmd.PersistentFlags().BoolVar(
		&options.ServeStatic, "serve-static", true, "Use project's router to serve static files")
	serveCmd.PersistentFlags().StringVar(
		&options.Secret, "secret", "s3cr3t", `Secret string to use for signing jwt.
 Note if you change it, already authenticated users would be logged out`)
	serveCmd.PersistentFlags().StringVar(
		&options.DbOptions.User, "dbuser", "postgres", "Database user")
	serveCmd.PersistentFlags().StringVar(
		&options.DbOptions.Password, "dbpassword", "", "Database password")
	serveCmd.PersistentFlags().StringVar(
		&options.DbOptions.Host, "dbhost", "localhost", "Database host")
	serveCmd.PersistentFlags().StringVar(
		&options.DbOptions.Port, "dbport", "5433", "Database port")
	serveCmd.PersistentFlags().StringVar(
		&options.DbOptions.DbName, "dbname", "postgres", "Database name")
	// TODO add dbtype
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start http server",
	Long: `Initialize database, mailer, bots, and routes to start answering http requests`,
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
		l.Log("msg", "Mock mailer is used")
		// db := app.NewMockDatabase()
		// l.Log("msg", "Mock database is used. Data is not persistent")
		db, err := app.NewGormDatabase(options.DbOptions)
		if err != nil {
			l.Log("panic", err)
			os.Exit(1)
		}
		l.Log("msg", "Connected to database")
		if options.ServeStatic {
			staticDir := filepath.Join(options.Root, options.StaticFiles)
			_, err = os.Stat(staticDir)
			if os.IsNotExist(err) {
				l.Log("panic", err)
			} else {
				l.Log("StaticDir", staticDir)
			}
		} else {
			options.StaticFiles = ""
			l.Log("msg", "Serving static files disabled")
		}
		if err != nil {
			l.Log("status", "exiting", "message", "Fix 'panic' errors above to serve http requests")
			os.Exit(1)
		}
		s := app.InitServer(l, t, m, db, options)
		// TODO start polling messages from bots
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			s.QuitCh <- struct{}{}
		}()
		s.ListenAndServe()
		l.Log("msg", "Bye!")
	},
}
