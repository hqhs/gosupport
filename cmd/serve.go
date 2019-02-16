package cmd

import (
	"strings"
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
	serveCmd.PersistentFlags().StringVar(
		&options.TgBotTokens, "tgtokens", "", `Quoted comma separated list of telegram
bot token to use for dashboard`)
	// TODO add dbtype
}

type checker struct {
	l kitlog.Logger
	err error
}

func (c *checker) Add(err error) {
	if err != nil {
		c.err = err
		c.l.Log("panic", err.Error())
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start http server",
	Long: `Initialize database, mailer, bots, and routes to start answering http requests`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// TODO use context for initialization process
		// TODO allow user set logger level (better to do it globally)
		l := kitlog.NewLogfmtLogger(os.Stdout)
		c := checker{l, nil}
		l.Log("root", options.Root, "address", options.Domain + ":" + options.Port)
		t, err := templator.NewTemplator(options.Root)
		c.Add(err)
		l.Log("templates", fmt.Sprintf("%+v", t.GetTemplates()))
		m := app.NewMockMailer(t, l)
		l.Log("msg", "Mock mailer is used")
		// db := app.NewMockDatabase()
		// l.Log("msg", "Mock database is used. Data is not persistent")
		db, err := app.InitPostgres(options.DbOptions)
		c.Add(err)
		l.Log("msg", "Connected to database")
		if options.ServeStatic {
			staticDir := filepath.Join(options.Root, options.StaticFiles)
			_, err = os.Stat(staticDir)
			if os.IsNotExist(err) {
				c.Add(err)
			} else {
				l.Log("StaticDir", staticDir)
			}
		} else {
			options.StaticFiles = ""
			l.Log("msg", "Serving static files disabled")
		}
		if c.err != nil {
			l.Log("status", "exiting", "message", "Fix 'panic' errors above to serve http requests")
			os.Exit(1)
		}
		s := app.InitServer(l, t, m, db, options)
		// TODO start polling messages from bots
		if len(options.TgBotTokens) > 0 {
			tok := strings.Split(options.TgBotTokens, ",")
			if len(tok) > 1 {
				l.Log("err", "Yet only one bot at a time is supported")
				os.Exit(1)
			}
			// Docker style management, since there's no guarantees what bot name
			// is unique, and I dont want to perform complex manipulations with
			// name/type, just generate 8 symbol md5 hashes of auth tokens and
			// store bots as map[hash]*Bot. It's better to use Bot interface for
			// state management only, therefore separate Connector type with
			// three channels : receive-only, write-only, and errors. Receiver
			// returns new messages from bot, which are broadcasted to socket hubs
			// later, and writer sends messages to chats with customers. Then
			// something goes wrong on either side, send err, try to notificate
			// user about it if hub is working and log it.
			t, err := app.NewTgBot(s, tok[0])
			if err != nil {
				l.Log("err", err, "then", "during initializing new telegram bot")
				os.Exit(1)
			}
			s.Add(t)
		}
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			// notify all bots to exit
			s.StopBots()
			// Quit server, server will block all incoming websocket messages,
			// but continue to send updates to dashboard from tg (then bots are
			// still quitting) and then bots are done, server will stop itself
			s.Shutdown()
		}()
		s.RunBots()
		s.ListenAndServe()
		l.Log("msg", "Bye!")
	},
}
