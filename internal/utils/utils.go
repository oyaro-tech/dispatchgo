package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ExpandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Failed to get user home directory: " + err.Error())
	}

	if path == "~" {
		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return "", errors.New("Failed to process file path: " + path)
}
