package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var base string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	base = filepath.Join(filepath.Dir(wd), "assets-storage")
	if err := os.MkdirAll(base, 0755); err != nil {
		panic(err)
	}
}

func SaveFile(file io.Reader, path string) (string, error) {
	cleanPath, err := SanitizePath(path)
	if err != nil {
		return "", err
	}

	abs := filepath.Join(base, cleanPath)

	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(abs)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(abs)
		return "", err
	}

	return path, nil
}

func SanitizePath(path string) (string, error) {
	path, err := filepath.Rel(base, filepath.Join(base, path))
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, "..") {
		return "", fmt.Errorf("invalid path")
	}

	return path, nil
}

func GetAbsPath(path string) string {
	return filepath.Join(base, path)
}
