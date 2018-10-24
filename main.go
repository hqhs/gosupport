package main

// Config describes every ENV variable used in project
type Config struct {
	Host       string
	MongoDBUrl string
	DBName     string
	Debug      bool
	Secret     string
}

func (c Config) String() {

}

// Conf represents settings from config.yaml
var Conf Config

func main() {
	ExecuteCmd()
}
