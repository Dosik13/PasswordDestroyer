package src

import (
	"go.uber.org/zap"
)

func NewLogger(debug bool) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	}
	return logger, nil

}
