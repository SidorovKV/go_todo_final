package config

import (
	"flag"
	"os"
)

type Config struct {
	Port   string
	DBFile string
}

func New() *Config {
	cfg := Config{}
	cfg.parseFlags()

	cfg.Port = ":" + cfg.Port

	port := os.Getenv("TODO_PORT")
	if len(port) > 0 {
		cfg.Port = ":" + port
	}

	dbfile := os.Getenv("TODO_DBFILE")
	if len(dbfile) > 0 {
		cfg.DBFile = dbfile
	}

	return &cfg
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Port, "httport", "7540", "")
	flag.StringVar(&c.DBFile, "dbfile", "./data/scheduler.db", "")

	flag.Parse()
}
