package main

import (
	"PasswordDestroyer/src"
	"fmt"
	"github.com/spf13/pflag"
)

func main() {
	debug := pflag.Bool("debug", false, "Enable debug mode")
	pflag.Parse()

	if len(pflag.Args()) < 2 {
		fmt.Println("Incorrect input! go run main.go <path_to_wordlist> <hash_to_crack> --debug")
		return
	}

	filePath := pflag.Arg(0)
	hash := pflag.Arg(1)

	logger, err := src.NewLogger(*debug)
	if err != nil {
		fmt.Println("Error initializing logger:", err)
	}

	defer logger.Sync()

	hasher := src.NewHasher(logger)

	hasher.Logger.Info("Starting the program")

	hasher.Run(filePath, hash)
}
