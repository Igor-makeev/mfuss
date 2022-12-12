package configs

type Config struct {
	SrvAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}
