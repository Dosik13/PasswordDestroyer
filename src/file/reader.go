package file

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

func GetAllPasswordsFromFile(filePath string) ([]string, error) {
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
