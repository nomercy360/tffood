package config

type Default struct {
	JWTSecret   string `env:"JWT_SECRET,required"`
	DBPath      string `env:"DB_PATH" envDefault:"./app.db"`
	BotToken    string `env:"TELEGRAM_BOT_TOKEN,required"`
	CdnURL      string `env:"CDN_URL,required"`
	ExternalURL string `env:"EXTERNAL_URL,required"`
	WebAppURL   string `env:"WEB_APP_URL,required"`
	AWS         AWSConfig
	Server      ServerConfig
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
}

type AWSConfig struct {
	AccessKey string `env:"AWS_ACCESS_KEY_ID,required"`
	SecretKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	Bucket    string `env:"AWS_BUCKET,required"`
	Endpoint  string `env:"AWS_ENDPOINT,required"`
}
