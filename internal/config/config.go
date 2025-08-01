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
	}

	Server struct {
		Host        string        `env:"SERVER_HOST"          env-default:"0.0.0.0"`
		Port        int           `env:"SERVER_PORT"          env-default:"8081"`
		Timeout     time.Duration `env:"SERVER_TIMEOUT"       env-default:"5s"`
		IdleTimeout time.Duration `env:"SERVER_IDDLE_TIMEOUT" env-default:"60s"`
	}

	Auth struct {
		AdminToken  string        `env:"AUTH_ADMIN_TOKEN" env-default:"my-secret-admin-token"`
		TokenSecret string        `env:"TOKEN_SECRET"     env-default:"my-secret-token"`
		TokenTTL    time.Duration `env:"TOKEN_TTL"        env-default:"1h"`
	}

	Cache struct {
		RedisHost  string `env:"REDIS_HOST"         env-default:"redis"`
		RedisPort  int    `env:"REDIS_PORT"         env-default:"6379"`
		TTLSeconds int    `env:"CACHE_TTL_SECONDS"  env-default:"300"`
		MaxItems   int    `env:"CACHE_MAX_ITEMS"    env-default:"1000"`
	}
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		log.Fatal("config path is empty")
	}

	if err := godotenv.Load(path); err != nil {
		log.Fatalf("cannot load env file: %s", err)
	}

	log.Println("using env file: " + path)

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("failed to read env: " + err.Error())
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
