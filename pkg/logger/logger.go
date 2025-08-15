package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logPath = "./logs/go.log"

func NewLogger(env string) *zap.Logger {
	var cfg zap.Config

	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "cannot create log dir: %v\n", err)
		os.Exit(1)
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open log file %q: %v\n", logPath, err)
		os.Exit(1)
	}
	f.Close()

	switch env {
	case "local":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//cfg.OutputPaths = []string{"stdout", logPath}
	cfg.OutputPaths = []string{logPath}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}
