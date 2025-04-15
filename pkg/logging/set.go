package logging

import (
	"github.com/sirupsen/logrus"

	"management/pkg/logging/hook"
	"management/pkg/logging/logcat"
	"management/pkg/logging/syslog"
	"management/pkg/logging/txt"
)

// SetTextFormatter set the text formatter for given logger.
func SetTextFormatter(logger *logrus.Logger) {
	logger.Formatter = txt.NewTextFormatter()
	logger.ReportCaller = true
	logger.AddHook(hook.NewContextHook())
}

// SetSyslogFormatter set the text formatter for given logger.
func SetSyslogFormatter(logger *logrus.Logger) {
	logger.Formatter = syslog.NewSyslogFormatter()
	logger.ReportCaller = true
	logger.AddHook(hook.NewContextHook())
}

// SetJSONFormatter set the JSON formatter for given logger.
func SetJSONFormatter(logger *logrus.Logger) {
	logger.Formatter = &logrus.JSONFormatter{}
	logger.ReportCaller = true
	logger.AddHook(hook.NewContextHook())
}

// SetLogcatFormatter set the logcat formatter for given logger.
func SetLogcatFormatter(logger *logrus.Logger) {
	logger.Formatter = logcat.NewLogcatFormatter()
	logger.ReportCaller = true
	logger.AddHook(hook.NewContextHook())
}
