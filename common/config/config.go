package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Database struct {
	Host         string `env:"DATABASE_HOST,default=localhost"`
	Port         string `env:"DATABASE_PORT,default=5432"`
	Username     string `env:"DATABASE_USERNAME,required"`
	Password     string `env:"DATABASE_PASSWORD,required"`
	Name         string `env:"DATABASE_NAME,required"`
	SSLMode      string `env:"DATABASE_SSL_MODE,default=disable"`
	MaxOpenConns int    `env:"DATABASE_MAX_OPEN_CONNS,default=5"`
	MaxIdleConns int    `env:"DATABASE_MAX_IDLE_CONNS,default=1"`
}

type Auth struct {
	PrivateKey string `env:"PRIVATE_KEY,required"`
	PublicKey  string `env:"PUBLIC_KEY,required"`
}

type Config struct {
	HashIDSalt      string `env:"HASHID_SALT"`
	HashIDMinLength int    `env:"HASHID_MIN_LENGTH"`
	Port            string `env:"PORT,default=6666"`
	Auth            Auth
	Database        Database
}

func NewConfig(env string) (*Config, error) {
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "[NewConfig] error decoding env")
	}

	return &config, nil
}
