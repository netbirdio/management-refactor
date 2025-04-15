package db

type config struct {
	Engine         string `env:"NB_STORE_ENGINE" envDefault:"sqlite"`
	PostgresDsnEnv string `env:"NB_STORE_ENGINE_POSTGRES_DSN" envDefault:""`
}
