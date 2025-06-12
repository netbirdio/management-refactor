package db

type Config struct {
	Engine      string `env:"NB_STORE_ENGINE" envDefault:"sqlite"`
	PostgresDsn string `env:"NB_STORE_ENGINE_POSTGRES_DSN" envDefault:""`
	DataDir     string `env:"NB_STORE_DATA_DIR" envDefault:"/var/lib/netbird"`
}
