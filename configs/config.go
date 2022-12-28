package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Config struct {
	SrvAddr         string `env:"SERVER_ADDRESS" `
	BaseURL         string `env:"BASE_URL" `
	FileStoragePath string `env:"FILE_STORAGE_PATH" `
	DBDSN           string `env:"DATABASE_DSN"`
}

func NewConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.SrvAddr, "a", "localhost:8080", "server addres to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "shortener base URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "file_storage.txt", "path to storage file")
	flag.StringVar(&cfg.DBDSN, "d", "", "database adress")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("failed to parse config environment variables")
	}

	logrus.Printf("env variable SERVER_ADDRESS=%v", cfg.SrvAddr)
	logrus.Printf("env variable BASE_URL=%v", cfg.BaseURL)
	logrus.Printf("env variable FILE_STORAGE_PATH=%v", cfg.FileStoragePath)

	return &cfg
}
