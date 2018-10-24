package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gosupport",
	Short: "Management commands for support bot dashboard",
	Long:  "From top to bottom re-written version of usefull tool for in-ompany helpdesking.",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(viper.GetString("host"))
		Runserver()
	},
}

var createSuperUser = &cobra.Command{

	Use:   "createsuperuser",
	Short: "Create admin with god's powers",
	Long:  "Connect to database, ask for credentials, and save them to Helpdesker model",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter email: ")
		email, _ := reader.ReadString('\n')
		sz := len(email)
		email = email[:sz-1] // delete newline character

		if email == "" {
			fmt.Printf("Email should not be empty")
			return
		}

		fmt.Print("Enter password: \n")
		password, _ := terminal.ReadPassword(0)

		fmt.Print("Confirm password: \n")
		confirmation, _ := terminal.ReadPassword(0)

		if string(password) != string(confirmation) {
			fmt.Printf("ERROR: Passwords don't match\n")
			return
		}

		h := Helpdesker{Email: email, IsAdmin: true, IsActive: true}
		h.SetPassword(string(password))

		collection := GetHelpdeskerCollection(Session)
		if err := h.Save(collection); err != nil {
			fmt.Printf("Couldn't create new superuser, error: %v\n", err)
			return
		}

		fmt.Printf("Superuser created successfuly, email: %v, password: %v\n", email, string(password))
	},
}

var checkDatabase = &cobra.Command{
	Use:   "checkdatabase",
	Short: "Create database connection and print usage stats",
	Long: `Database connection created in store.Init() method,
	      but it's easier to debug then separate command exists`,
	Run: func(cmd *cobra.Command, args []string) {
		if Session == nil {
			fmt.Printf("Couldn't connect to database, you should see error log above\n")
		}
		fmt.Printf("Database connection successfuly created\n")
	},
}

// ExecuteCmd executes given command and exit if recieved en error
func ExecuteCmd() {
	// NOTE here should be all initializations, which should be initialized for every command
	// Since database used in all commands, I initialize it here
	var err error
	Session, err = mgo.Dial(Conf.MongoDBUrl)
	defer Session.Close()

	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	}
	// Ensure what unique values is actually unique
	err = GetHelpdeskerCollection(Session).EnsureIndex(index)
	if err != nil {
		log.Printf("Couldn't ensure index for helpdeskers: %v\n", err)
		panic(err)
	}
	log.Printf("Database index ensured\n")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// initialize config
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(createSuperUser)
	rootCmd.AddCommand(checkDatabase)
	// add flags
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	Must(viper.ReadInConfig())

	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())

	// read config values from config.yaml
	Must(viper.Unmarshal(&Conf))

	fmt.Printf("Configuration values: %v\n", Conf)
}
