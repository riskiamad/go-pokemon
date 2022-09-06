package configs

import (
	"os"
)

var Config = loadConfig()

type config struct {
	Host     string
	Port     string
	MongoURI string
	DBName   string
}

func loadConfig() *config {
	c := new(config)

	c.Host = os.Getenv("HOST")
	c.Port = os.Getenv("PORT")
	c.MongoURI = os.Getenv("MONGOURI")
	c.DBName = os.Getenv("DBNAME")

	return c
}
