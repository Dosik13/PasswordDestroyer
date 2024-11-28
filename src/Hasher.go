package src

import (
	"bufio"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"sync"
)

type Hasher struct {
	Passwords []string
	Logger    *zap.Logger
}

func NewHasher(logger *zap.Logger) *Hasher {
	return &Hasher{Logger: logger}
}

func (h *Hasher) GetAllPasswordsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var passwords []string
	reader := bufio.NewReader(file)

	for {
		password, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			passwords = append(passwords, strings.TrimSpace(password))
			break
		}
		if err != nil {
			return nil, err
		}

		passwords = append(passwords, strings.TrimSpace(password))
	}
	return passwords, nil
}

func (h *Hasher) findHash(done chan struct{}, hash, password string, wg *sync.WaitGroup, found chan<- string, isMD5 bool, logger *zap.Logger) {
	defer wg.Done()

	select {
	case <-done:
		logger.Debug("Received done signal, stopping worker")
		return
	default:
		logger.Debug("Checking password", zap.String("password", password))
		if checkHash(hash, password, isMD5) {
			found <- password
			return
		}
	}
}

func (h *Hasher) StartWorkers(done chan struct{}, hash string, found chan<- string, isMD5 bool, wg *sync.WaitGroup) {
	h.Logger.Info("Starting workers", zap.Int("goroutines", len(h.Passwords)))

	for _, password := range h.Passwords {
		wg.Add(1)
		go h.findHash(done, hash, password, wg, found, isMD5, h.Logger)
	}
}

func (h *Hasher) Run(filePath, hash string) {
	passwrds, _ := h.GetAllPasswordsFromFile(filePath)
	h.Passwords = passwrds

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file!")
		return
	}
	defer f.Close()

	var wg sync.WaitGroup

	found := make(chan string, 1)
	done := make(chan struct{})

	isMD5 := len(hash) == 32

	h.StartWorkers(done, hash, found, isMD5, &wg)

	isFound := false

	select {
	case password := <-found:
		h.Logger.Info("Match found!", zap.String("password", password))
		isFound = true
		close(done)
	case <-done:
		h.Logger.Debug("Workers are being cancelled")
	}

	wg.Wait()
	if !isFound {
		h.Logger.Info("No matching hash found!")
	}
}

func checkHash(hash, password string, isMD5 bool) bool {
	if isMD5 {
		return hash == ToMD5(password)
	}
	return hash == ToSha256(password) || hash == ToKeccak256(password)
}
