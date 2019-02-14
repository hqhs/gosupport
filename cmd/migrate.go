package cmd

import (
	"os"

	"github.com/hqhs/gosupport/internal/app"
	"github.com/spf13/cobra"
	kitlog "github.com/go-kit/kit/log"
)

func init() {
	migrateCmd.PersistentFlags().StringVar(
		&options.DbOptions.User, "dbuser", "postgres", "Database user")
	migrateCmd.PersistentFlags().StringVar(
		&options.DbOptions.Password, "dbpassword", "", "Database password")
	migrateCmd.PersistentFlags().StringVar(
		&options.DbOptions.Host, "dbhost", "localhost", "Database host")
	migrateCmd.PersistentFlags().StringVar(
		&options.DbOptions.Port, "dbport", "5433", "Database port")
	migrateCmd.PersistentFlags().StringVar(
		&options.DbOptions.DbName, "dbname", "postgres", "Database name")
	// TODO add dbtype
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Use gorm auto migration feature for database schema",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// TODO allow user set logger level (better to do it globally)
		l := kitlog.NewLogfmtLogger(os.Stdout)
		l.Log("root", options.Root, "address", options.Domain + ":" + options.Port)
		gormDB, err := app.NewGormDatabase(options.DbOptions)
		if err != nil {
			l.Log("panic", err)
			l.Log("status", "exiting", "message", "Fix 'panic' errors above to serve http requests")
			os.Exit(1)
		}
		gormDB.AutoMigrate(&app.Admin{}, &app.User{}, &app.Message{})
	},
}
