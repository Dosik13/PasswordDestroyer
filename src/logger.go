package src

import (
	"go.uber.org/zap"
)

func NewLogger(debug bool) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
