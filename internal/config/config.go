package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr              string
	DatabaseURI          string
	AccrualSystemAddress string
}

func ParseFlags() (Config, error) {
	cfg := Config{
		RunAddr: ":8080",
		AccrualSystemAddress: "http://localhost:8081",
	}

	if envRunAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.RunAddr = envRunAddr
	}

	if envDatabaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.DatabaseURI = envDatabaseURI
	}

	if accrualAddress, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		cfg.AccrualSystemAddress = accrualAddress
	}

	flag.StringVar(&cfg.RunAddr, "a", cfg.RunAddr, "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database connection")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "address for connecting to accrual")

	flag.Parse()

	return cfg, nil
}
