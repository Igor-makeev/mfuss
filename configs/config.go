// Пакет конфигурации приложения
package configs

import (
	"crypto/tls"
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config - тип структуры конфигурации приложения.
type Config struct {
	SrvAddr         string `env:"SERVER_ADDRESS" `
	BaseURL         string `env:"BASE_URL" `
	FileStoragePath string `env:"FILE_STORAGE_PATH" `
	DBDSN           string `env:"DATABASE_DSN"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
	TLSConf         *tls.Config
}

// Конструктор конфигов
func NewConfig() *Config {
	// Инизциализируем конфиг.
	var cfg Config

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	cfg.SrvAddr = viper.GetString("conf.server_address")
	cfg.BaseURL = viper.GetString("conf.base_url")
	cfg.FileStoragePath = viper.GetString("conf.file_storage_path")
	cfg.DBDSN = viper.GetString("conf.database_dsn")
	cfg.EnableHTTPS = viper.GetBool("conf.enable_https")

	flag.StringVar(&cfg.SrvAddr, "a", ":8080", "server addres to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "shortener base URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "file_storage.txt", "path to storage file")
	flag.StringVar(&cfg.DBDSN, "d", "", "database adress")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "https")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("failed to parse config environment variables")
	}
	if cfg.EnableHTTPS {

		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem") //загрузка серверного сертификата и ключа
		if err != nil {
			log.Fatal(err)
		}

		cfg.TLSConf = &tls.Config{Certificates: []tls.Certificate{cert}}
	}
	logrus.Printf("env variable SERVER_ADDRESS=%v", cfg.SrvAddr)
	logrus.Printf("env variable BASE_URL=%v", cfg.BaseURL)
	logrus.Printf("env variable FILE_STORAGE_PATH=%v", cfg.FileStoragePath)
	logrus.Printf("env variable DATABASE_DSN=%v", cfg.DBDSN)
	logrus.Printf("env variable ENABLE_HTTPS=%v", cfg.EnableHTTPS)

	return &cfg
}
