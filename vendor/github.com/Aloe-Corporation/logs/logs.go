package logs

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DEBUG is an identifier for debug level.
	DEBUG = "DEBUG"
	// INFO is an identifier for information level.
	INFO = "INFO"
	// WARN is an identifier for warning level.
	WARN = "WARN"
	// ERROR is an identifier for error level.
	ERROR = "ERROR"
)

var (
	Config Conf
	logger zap.Logger
)

// Conf for logs package.
type Conf struct {
	// Level allowed values : DEBUG, INFO, WARN, ERROR
	Level  string   `yaml:"level"`
	Output []string `yaml:"output"`
}

func init() {
	// Set default logger
	var err error
	defaultConfLogger := Conf{
		Level:  DEBUG,
		Output: []string{"stdout"},
	}
	logger, err = FactoryLogger(defaultConfLogger)
	if err != nil {
		fmt.Println(fmt.Errorf("fail to init default logger: %w", err))
	} else {
		logger.Info("logger construction succeeded", zap.String("service", "logger"))
	}
}

// Init the logger with Config.
func Init() error {
	var err error
	logger, err = FactoryLogger(Config)
	return err
}

// Get the instance of the zap.Logger.
func Get() *zap.Logger {
	return &logger
}

// FactoryLogger instanciate a new logger with conf (default level is DEBUG).
func FactoryLogger(c Conf) (zap.Logger, error) {
	lvl := zapcore.DebugLevel
	switch c.Level {
	case INFO:
		lvl = zapcore.InfoLevel
	case WARN:
		lvl = zapcore.WarnLevel
	case ERROR:
		lvl = zapcore.ErrorLevel
	}

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(lvl),
		OutputPaths:      c.Output,
		ErrorOutputPaths: c.Output,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	l, err := config.Build()
	if err != nil {
		return *l, fmt.Errorf("fail to build logger: %w", err)
	}

	return *l, nil
}
