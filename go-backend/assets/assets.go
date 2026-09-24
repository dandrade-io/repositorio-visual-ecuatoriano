package assets

import (
	"errors"
	"io"
	"net/http"
)

func mimeType(file io.Reader) (string, error) {
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return http.DetectContentType(buf), nil
}
