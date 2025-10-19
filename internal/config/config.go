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
	cfg := Config{
		RunAddr: ":8080",
	}

	if envRunAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.RunAddr = envRunAddr
	}

	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.DatabaseURI = envDatabaseDSN
	}

	flag.StringVar(&cfg.RunAddr, "a", cfg.RunAddr, "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database connection")

	flag.Parse()

	return cfg, nil
}
