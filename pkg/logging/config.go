package logging

type LoggingConfig struct {
	LogLevels map[string]string `mapstructure:"log_levels"`
}
