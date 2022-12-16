package configs

type Config struct {
	SrvAddr         string `env:"SERVER_ADDRESS" `
	BaseURL         string `env:"BASE_URL" `
	FileStoragePath string `env:"FILE_STORAGE_PATH" `
}
