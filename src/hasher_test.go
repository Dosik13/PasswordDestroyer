package src

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"testing"
)

func newTestLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	return logger
}

func TestFindHash(t *testing.T) {
	logger := newTestLogger()
	hasher := NewHasher(logger)

	tests := []struct {
		name     string
		hash     string
		password string
		hashType HashType
		expected string
	}{
		{
			name:     "SHA256 match",
			hash:     "fbc9e0b78a38aa356a29f1ee49f43ef045a8a4912f483f79b16d122b4b9ab2ea",
			password: "amublance34",
			hashType: Other,
			expected: "amublance34",
		},
		{
			name:     "SHA256 match",
			hash:     "532f011ec89ff0e2e1be76953593b588d47e8a454d18f554e4f2ea6d89615a10",
			password: "dolphin",
			hashType: Other,
			expected: "dolphin",
		},
		{
			name:     "MD5 match",
			hash:     "36cdf8b887a5cffc78dcd5c08991b993",
			password: "dolphin",
			hashType: MD5,
			expected: "dolphin",
		},
		{
			name:     "No match",
			hash:     "fbc9e0b78a38aa356a29f1ee49f43ef045a8a4912f483f79b16d122b4b9ab2ea",
			password: "dolphin",
			hashType: MD5,
			expected: "",
		},
		{
			name:     "Done signal",
			hash:     "fbc9e0b78a38aa356a29f1ee49f43ef045a8a4912f483f79b16d122b4b9ab2ea",
			password: "dolphin",
			hashType: MD5,
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			done := make(chan struct{})
			found := make(chan string, 1)
			var wg sync.WaitGroup

			if test.name == "Done signal" {
				close(done)
			}

			wg.Add(1)
			go hasher.findHash(done, test.hash, test.password, &wg, found, test.hashType, logger)
			wg.Wait()

			select {
			case password := <-found:
				if password != test.expected {
					t.Errorf("Expected password '%s', got '%s'", test.expected, password)
				}
			default:
				if test.expected != "" {
					t.Errorf("Expected to find password, but found channel is empty")
				}
			}
		})
	}
}

func TestStartWorkers(t *testing.T) {
	logger := newTestLogger()
	hasher := NewHasher(logger)
	hasher.(*Hasher).Passwords = []string{"pink", "dolphin", "amublance34", "Pelican123"}
	done := make(chan struct{})
	found := make(chan string, 1)
	var wg sync.WaitGroup

	hasher.startWorkers(done, "fbc9e0b78a38aa356a29f1ee49f43ef045a8a4912f483f79b16d122b4b9ab2ea", found, Other, &wg)

	wg.Wait()

	select {
	case password := <-found:
		if password != "amublance34" {
			t.Errorf("Expected password 'amublance34', got '%s'", password)
		}
	default:
		t.Error("Expected to find password, but found channel is empty")
	}
}

func TestHasherRun(t *testing.T) {
	// Set up a shared logger
	logger := newTestLogger()

	hasher := NewHasher(logger)

	hasher.getAllPasswordsFromFile("passwords_test.txt")

	tests := []struct {
		name       string
		filePath   string
		hash       string
		expected   bool
		shouldFail bool
	}{
		{
			name:       "File not found error",
			filePath:   "file123.txt",
			hash:       "fbc9e0b78a38aa356a29f1ee49f43ef045a8a4912f483f79b16d122b4b9ab2ea",
			expected:   false,
			shouldFail: true,
		},
		{
			name:       "Hash found",
			filePath:   "passwords_test.txt",
			hash:       "68ad1d70106186e2e8a05ee8352a1ec9733ae75f87f255c79bca15d6afe79a7f",
			expected:   true,
			shouldFail: false,
		},
		{
			name:       "Hash not found",
			filePath:   "passwords_test.txt",
			hash:       "fbc9e0b78a38aa356a29f1ee49f43ef045a8a4912f483f79b16d122b4b9ab2ea",
			expected:   false,
			shouldFail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			found, err := hasher.Run(tc.filePath, tc.hash)

			if tc.shouldFail && err == nil {
				t.Fatalf("Expected an error but got nil")
			} else if !tc.shouldFail && err != nil {
				t.Fatalf("Did not expect an error but got: %v", err)
			}

			if found != tc.expected {
				t.Errorf("Expected found = %v, got %v", tc.expected, found)
			}
		})
	}
}
