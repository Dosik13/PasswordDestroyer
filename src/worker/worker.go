package worker

import (
	"PasswordDestroyer/src/hashing"
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
)

func findHash(ctx context.Context, hash string, passwords <-chan string, wg *sync.WaitGroup, res chan<- string, isMD5 bool, logger *zap.Logger) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Debug("Context canceled, stopping worker.")
			return
		case password, ok := <-passwords:
			if !ok {
				logger.Debug("Passwords channel closed, stopping worker.")
				return
			}
			logger.Debug("Checking password", zap.String("password", password))
			if checkHash(hash, password, isMD5) {
				logger.Info("Password match found", zap.String("password", password))
				res <- password
				return
			}
		}
	}
}

func StartWorkers(ctx context.Context, hash string, words <-chan string, res chan<- string, isMD5 bool, wg *sync.WaitGroup, logger *zap.Logger) {
	goroutines := 200
	logger.Info("Starting workers", zap.Int("goroutines", goroutines))

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		logger.Debug("Starting worker", zap.Int("worker_id", i))
		go findHash(ctx, hash, words, wg, res, isMD5, logger)
	}
}

func SendPasswordsToChannel(passwords []string, words chan<- string, logger *zap.Logger) {
	go func() {
		for _, password := range passwords {
			logger.Debug("Sending password to channel", zap.String("password", password))
			words <- password
		}
		close(words)
		logger.Info("Password channel closed")
	}()
}

func HandleResult(ctx context.Context, res <-chan string, cancel context.CancelFunc, wg *sync.WaitGroup, logger *zap.Logger) {
	found := false
	go func() {
		select {
		case password := <-res:
			logger.Info("Match found!", zap.String("password", password))
			fmt.Printf("Match found! The password is %s\n", password)
			cancel()
			found = true
		case <-ctx.Done():
			logger.Debug("Context canceled")
		}
	}()

	wg.Wait()
	if !found {
		logger.Info("No matching hash found!")
		fmt.Println("No matching hash found!")
	}
}

func checkHash(hash, password string, isMD5 bool) bool {
	if isMD5 {
		return hash == hashing.ToMD5(password)
	}
	return hash == hashing.ToSha256(password) || hash == hashing.ToKeccak256(password)
}
