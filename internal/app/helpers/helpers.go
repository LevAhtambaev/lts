package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
)

func LoadImage(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("ошибка при получении информации о файле: %w", err)
	}

	size := fileInfo.Size()
	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	base64Encoding := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buffer)

	return base64Encoding, nil
}
