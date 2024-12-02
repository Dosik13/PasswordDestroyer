package src

import (
	"bufio"
	"errors"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"time"
)

// const ErrorNotFound := errors.New("Not found")

type HashType int

const (
	Other HashType = iota
	MD5
)

const Time = 500 * time.Millisecond

type Hasherer interface {
	getAllPasswordsFromFile(filePath string) error
	findHash(done chan struct{}, hash, password string, found chan<- string, hashType HashType, logger *zap.Logger)
	startWorkers(done chan struct{}, hash string, found chan<- string, hashType HashType)
	Run(filePath, hash string) (bool, error)
}

type Hasher struct {
	Passwords []string
	Logger    *zap.Logger
}

func NewHasher(logger *zap.Logger) Hasherer {
	return &Hasher{Logger: logger}
}

func (h *Hasher) getAllPasswordsFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		password, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			h.Passwords = append(h.Passwords, strings.TrimSpace(password))
			break
		}
		if err != nil {
			return err
		}

		h.Passwords = append(h.Passwords, strings.TrimSpace(password))
	}
	return nil
}

func (h *Hasher) findHash(done chan struct{}, hash, password string, found chan<- string, hashType HashType, logger *zap.Logger) {
	select {
	case <-done:
		logger.Debug("Received done signal, stopping worker")
		return
	default:
		logger.Debug("Checking password", zap.String("password", password))
		if CheckHash(hash, password, hashType) {
			found <- password
			return
		}
	}
}

func (h *Hasher) startWorkers(done chan struct{}, hash string, found chan<- string, hashType HashType) {
	h.Logger.Info("Starting workers", zap.Int("goroutines", len(h.Passwords)))

	for _, password := range h.Passwords {
		go h.findHash(done, hash, password, found, hashType, h.Logger)
	}
}

func (h *Hasher) Run(filePath, hash string) (bool, error) {
	h.Logger.Info("Starting hash cracking", zap.String("hash", hash), zap.String("file", filePath))

	err := h.getAllPasswordsFromFile(filePath)
	if err != nil {
		return false, err
	}

	ticker := time.NewTicker(Time)
	found := make(chan string)
	done := make(chan struct{})

	var hashType HashType

	if len(hash) == 32 {
		hashType = MD5
	} else {
		hashType = Other
	}

	h.startWorkers(done, hash, found, hashType)

	isFound := false

	select {
	case password := <-found:
		h.Logger.Info("Match found!", zap.String("password", password))
		isFound = true
		close(done)
	case <-ticker.C:
		h.Logger.Debug("Ticker expired, stopping workers")
		close(done)
	}

	if !isFound {
		return false, nil
	}

	return true, nil
}
