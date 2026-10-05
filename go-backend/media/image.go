package media

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DetectMIME(r io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read for detection: %w", err)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek: %w", err)
	}

	return http.DetectContentType(buf[:n]), nil
}

func ToDataURI(r io.ReadSeeker) (string, error) {
	mimeType, err := DetectMIME(r)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
