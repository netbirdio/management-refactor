package activity

type config struct {
	Enabled bool `env:"NB_EVENT_ACTIVITY_LOG_ENABLED" envDefault:"true"`
}
