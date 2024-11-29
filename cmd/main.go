package main

import (
	"PasswordDestroyer/src"
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	debug := pflag.Bool("debug", false, "Enable debug mode")
	pflag.Parse()

	logger, err := src.NewLogger(*debug)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
	}

	defer logger.Sync()

	if len(pflag.Args()) < 2 {
		fmt.Println("Incorrect input! go run main.go <path_to_wordlist> <hash_to_crack> --debug")
		return
	}

	filePath := pflag.Arg(0)
	hash := pflag.Arg(1)

	hasher := src.NewHasher(logger)

	found, err := hasher.Run(filePath, hash)
	if err != nil {
		logger.Error("Error running hasher", zap.Error(err))
	}
	if !found {
		logger.Info("Password not found")
	}
}
