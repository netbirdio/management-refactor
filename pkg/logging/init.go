package logging

import (
	"fmt"
	"io"
	"log/syslog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const defaultLogSize = 5

// global map of package paths to *logrus.Logger
var (
	mu      sync.RWMutex
	loggers = map[string]*logrus.Logger{}
)

const basePrefix = "github.com/netbirdio/"

// Init reads the logging config from a YAML file and sets up per-package loggers.
func Init(configFilePath string) error {
	v := viper.New()
	v.SetConfigFile(configFilePath)

	// Read the YAML
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read logging config: %w", err)
	}

	var cfg LoggingConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal logging config: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()

	// For each package in our config, create a logger at the given level
	for pkgPath, levelStr := range cfg.LogLevels {
		l := logrus.New() // each package gets its own *logrus.Logger
		l.SetLevel(parseLogrusLevel(levelStr))

		// Optionally, set the formatter, output, etc.:
		// l.SetFormatter(&logrus.JSONFormatter{})
		// l.SetOutput(os.Stdout)

		loggers[pkgPath] = l
	}

	// Optionally, define a default logger for packages not explicitly listed
	if _, ok := loggers["default"]; !ok {
		defaultLogger := logrus.New()
		defaultLogger.SetLevel(logrus.InfoLevel)
		loggers["default"] = defaultLogger
	}

	return nil
}

// parseLogrusLevel is a helper that converts a string (e.g. "debug") to a logrus.Level.
func parseLogrusLevel(levelStr string) logrus.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	}
	// default
	return logrus.InfoLevel
}

// LoggerFor returns a *logrus.Logger for the specified package path.
// If there's no explicit logger, we return the "default" logger.
func LoggerFor(pkgPath string) *logrus.Logger {
	mu.RLock()
	defer mu.RUnlock()

	if l, ok := loggers[pkgPath]; ok {
		return l
	}
	if l, ok := loggers["default"]; ok {
		return l
	}

	logrus.Tracef("No logger configured for %q; using fallback (info-level) logger", pkgPath)
	fallback := logrus.New()
	fallback.SetLevel(logrus.InfoLevel)
	return fallback
}

func LoggerForThisPackage() *logrus.Logger {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return LoggerFor("default")
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return LoggerFor("default")
	}

	fullFuncName := fn.Name()
	pkgPath := parsePackageFromFuncName(fullFuncName)

	if strings.HasPrefix(pkgPath, basePrefix) {
		pkgPath = strings.TrimPrefix(pkgPath, basePrefix)
	}

	return LoggerFor(pkgPath)
}

func parsePackageFromFuncName(funcName string) string {
	parts := strings.Split(funcName, "/")
	if len(parts) == 0 {
		return "default"
	}
	last := parts[len(parts)-1]

	base := strings.Join(parts[:len(parts)-1], "/")

	dotIdx := strings.IndexByte(last, '.')
	var pkgName string
	if dotIdx == -1 {
		pkgName = last
	} else {
		pkgName = last[:dotIdx]
	}

	return base + "/" + pkgName
}

// InitLog parses and sets log-level input
func InitLog(logLevel string, logPath string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Errorf("Failed parsing log-level %s: %s", logLevel, err)
		return err
	}
	customOutputs := []string{"console", "syslog"}

	if logPath != "" && !slices.Contains(customOutputs, logPath) {
		maxLogSize := getLogMaxSize()
		lumberjackLogger := &lumberjack.Logger{
			// Log file absolute path, os agnostic
			Filename:   filepath.ToSlash(logPath),
			MaxSize:    maxLogSize, // MB
			MaxBackups: 10,
			MaxAge:     30, // days
			Compress:   true,
		}
		logrus.SetOutput(io.Writer(lumberjackLogger))
	} else if logPath == "syslog" {
		addSyslogHook()
	}

	//nolint:gocritic
	if os.Getenv("NB_LOG_FORMAT") == "json" {
		SetJSONFormatter(logrus.StandardLogger())
	} else if logPath == "syslog" {
		SetSyslogFormatter(logrus.StandardLogger())
	} else {
		SetTextFormatter(logrus.StandardLogger())
	}
	logrus.SetLevel(level)

	setGRPCLibLogger()

	return nil
}

func setGRPCLibLogger() {
	logOut := logrus.StandardLogger().Writer()
	if os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL") != "info" {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(io.Discard, logOut, logOut))
		return
	}

	var v int
	vLevel := os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	}

	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(logOut, logOut, logOut, v))
}

func getLogMaxSize() int {
	if sizeVar, ok := os.LookupEnv("NB_LOG_MAX_SIZE_MB"); ok {
		size, err := strconv.ParseInt(sizeVar, 10, 64)
		if err != nil {
			logrus.Errorf("Failed parsing log-size %s: %s. Should be just an integer", sizeVar, err)
			return defaultLogSize
		}

		logrus.Infof("Setting log file max size to %d MB", size)

		return int(size)
	}
	return defaultLogSize
}

func addSyslogHook() {
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err != nil {
		logrus.Errorf("Failed creating syslog hook: %s", err)
	}
	logrus.AddHook(hook)
}
