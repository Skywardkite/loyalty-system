package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr     string
	DatabaseURI string
}

func ParseFlags() (Config, error) {
	var cfg Config

	flag.StringVar(&cfg.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database connection")

	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.RunAddr = envRunAddr
	}

	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.DatabaseURI = envDatabaseDSN
	}

	return cfg, nil
}
