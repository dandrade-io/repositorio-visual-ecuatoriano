package media

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
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

// ToDataURI encodes r to Base64
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

func CompressImg(r io.Reader, quality int) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}
