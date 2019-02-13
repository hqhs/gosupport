package cmd

import (
	"os"

	"github.com/hqhs/gosupport/internal"
	"github.com/spf13/cobra"
	kitlog "github.com/go-kit/kit/log"
)


var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Simple tech-support dashboard for in-company usage",
	Long: `Bare-bones dashboard for managing chats with anyone, who needs help,
 with telegrambot interface`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// TODO allow user set logger level globally
		logger := kitlog.NewLogfmtLogger(os.Stdout)
		logger.Log("root", options.Root, "address", options.Domain + ":" + options.Port)
		templator, err := internal.NewTemplator(options.Root)
		if err != nil {
			// NOTE we dont panic here to allow init process to finish and find all
			// errors, since they're independent. But http router serving won't start.
			logger.Log(err)
		}
		// m, err := newMockMailer()
		// if err != nil {
		// 	panic(err)
		// }
		//
		//
		// Legacy code:
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
		s := internal.InitServer(logger, templator, options)
		s.ServeRouter()
	},
}
