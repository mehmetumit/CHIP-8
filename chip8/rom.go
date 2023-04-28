package chip8

import (
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	return data, err
}
