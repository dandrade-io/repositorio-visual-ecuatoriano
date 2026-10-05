package albumcover

import (
	"encoding/json"
	"net/http"
	"repository/media"
	transporthttp "repository/transport/http"
	"strings"
)

type Handler struct{}

func (h *Handler) Metadata(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: err.Error(), Cause: err}
	}
	defer r.MultipartForm.RemoveAll()

	headers := r.MultipartForm.File["asset"]
	if len(headers) != 1 {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "only one asset must be provided"}
	}

	header := headers[0]
	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	mime, err := media.DetectMIME(file)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(mime, "image/") {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "only images are accepted"}
	}

	img, err := media.ToDataURI(file)
	if err != nil {
		return err
	}

	metadata, err := Extract(r.Context(), img)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(&metadata); err != nil {
		return err
	}

	return nil
}
