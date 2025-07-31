package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	MigrationDir string `env:"MIGRATION_DIR" env-default:"./migrations"`

	DB struct {
		Host     string `env:"DB_HOST"     env-default:"postgres"`
		Port     int    `env:"DB_PORT"     env-default:"5432"`
		User     string `env:"DB_USER"     env-default:"postgres"`
		Password string `env:"DB_PASSWORD" env-default:"postgres"`
		Name     string `env:"DB_NAME"     env-default:"documents"`
	} `env-prefix:"DB_"`

	Server struct {
		Host        string        `env:"SERVER_HOST"          env-default:"0.0.0.0"`
		Port        int           `env:"SERVER_PORT"          env-default:"8081"`
		Timeout     time.Duration `env:"SERVER_TIMEOUT"       env-default:"5s"`
		IdleTimeout time.Duration `env:"SERVER_IDDLE_TIMEOUT" env-default:"60s"`
	} `env-prefix:"SERVER_"`

	Auth struct {
		AdminToken string        `env:"AUTH_ADMIN_TOKEN" env-default:"my-secret-admin-token"`
		TokenTTL   time.Duration `env:"TOKEN_TTL"        env-default:"1h"`
	} `env-prefix:"AUTH_"`

	Cache struct {
		RedisHost  string `env:"REDIS_HOST"         env-default:"redis"`
		RedisPort  int    `env:"REDIS_PORT"         env-default:"6379"`
		TTLSeconds int    `env:"CACHE_TTL_SECONDS"  env-default:"300"`
		MaxItems   int    `env:"CACHE_MAX_ITEMS"    env-default:"1000"`
	} `env-prefix:"CACHE_"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		log.Fatal("config path is empty")
	} else {
		if err := godotenv.Load(path); err != nil {
			log.Fatalf("cannot load env file: %s", err)
		}
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist: " + configPath)
	}
	log.Println("using config file: " + configPath)

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
