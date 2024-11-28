package main

import (
	"PasswordDestroyer/src/file"
	"PasswordDestroyer/src/worker"
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"os"
	"sync"
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

	var logger *zap.Logger
	if *debug {
		var err error
		logger, err = zap.NewDevelopment()
		fmt.Println(" Ima li?")
		if err != nil {
			fmt.Println("Error initializing logger:", err)
		}
	} else {
		logger = zap.NewNop()
	}
	defer logger.Sync()

	logger.Info("Starting the program")

	passwords, _ := file.GetAllPasswordsFromFile(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file!")
		return
	}
	defer f.Close()

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	words := make(chan string, len(passwords))
	res := make(chan string, 1)

	isMD5 := len(hash) == 32

	worker.StartWorkers(ctx, hash, words, res, isMD5, &wg, logger)
	worker.SendPasswordsToChannel(passwords, words, logger)
	worker.HandleResult(ctx, res, cancel, &wg, logger)
}
