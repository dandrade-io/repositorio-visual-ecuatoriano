package assets

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"repository/database"
	"repository/embedding"
	"repository/files"
	"repository/media"
	transporthttp "repository/transport/http"
	"strconv"
	"strings"
	"uuid"
)

type Handler struct {
	Svc Service
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) error {
	claim, _ := transporthttp.ClaimsFrom(r.Context())
	userID, _ := uuid.Parse(claim.Subject)

	embed := false
	if val, err := strconv.ParseBool(r.URL.Query().Get("embed")); err == nil {
		embed = val
	}

	info, err := parseMultipartFile(r, embed)
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		return err
	}
	defer info.file.Close()

	asset, err := h.Svc.Create(r.Context(), CreateInput{File: info.file, Kind: info.kind, UserID: userID, Embed: embed})
	if err != nil {
		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to store asset", Cause: err}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&assetResp{asset.ID, asset.Path}); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Embed(w http.ResponseWriter, r *http.Request) error {
	req := &embedReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: err.Error()}
	}

	model := embedding.Model(req.Model)
	if !model.Valid() {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model"}
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "invalid asset id"}
	}

	if err := h.Svc.EmbedByAssetID(r.Context(), id, model); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "invalid asset id", Cause: err}
		}

		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to extract embedding", Cause: err}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	path, err := files.SanitizePath(r.PathValue("path"))
	if err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "invalid path"}
	}

	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, files.GetAbsPath(path))
	return nil
}

func (h *Handler) SearchByText(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("q")
	kind := r.URL.Query().Get("kind")
	model := embedding.JinaClipV2
	limit := 10

	if strings.TrimSpace(q) == "" {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "a query value must be provided"}
	}
	if strings.TrimSpace(kind) == "" {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "a kind value must be provided"}
	}
	if kind != "image" && kind != "audio" {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "kind value must be \"image\" or \"audio\" "}
	}
	if name := strings.TrimSpace(r.URL.Query().Get("model")); name != "" {
		model = embedding.Model(name)
		if !model.Valid() {
			return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model provided"}
		}
	}
	if num, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = num
	}

	result, err := h.Svc.SearchByText(r.Context(), q, model, limit)
	if err != nil {
		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to search for k-best candidates", Cause: err}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&result); err != nil {
		return err
	}

	return nil
}

func (h *Handler) SearchByAsset(w http.ResponseWriter, r *http.Request) error {
	modality := embedding.Modality(r.URL.Query().Get("modality"))
	model := embedding.Model(r.URL.Query().Get("model"))
	limit := 10

	if !modality.Valid() {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown modality provided, it must be \"text\", \"image\" or \"audio\""}
	}
	if !model.Valid() {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model provided"}
	}
	if num, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = num
	}

	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "not able to parse multipart form", Cause: err}
	}

	headers := r.MultipartForm.File["file"]
	if len(headers) != 1 {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "only one asset must be provided"}
	}

	file, err := headers[0].Open()
	if err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unable to process file", Cause: err}
	}
	defer file.Close()

	result, err := h.Svc.SearchByAsset(r.Context(), file, model, modality, limit)
	if err != nil {
		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to search for k-best candidates", Cause: err}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&result); err != nil {
		return err
	}

	return nil
}

type multipartInfo struct {
	file  multipart.File
	model embedding.Model
	kind  string
}

func parseMultipartFile(r *http.Request, embed bool) (multipartInfo, error) {
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: err.Error(), Cause: err}
	}

	headers := r.MultipartForm.File["file"]
	models := r.MultipartForm.Value["model"]
	if len(headers) != 1 {
		return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "only one asset must be provided"}
	}

	var model embedding.Model
	if embed {
		if len(models) != 1 {
			return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "only one model value must be provided"}
		}

		model = embedding.Model(models[0])
		if !model.Valid() {
			return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model provided"}
		}
	}

	header := headers[0]
	file, err := header.Open()
	if err != nil {
		return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unable to process file", Cause: err}
	}

	mime, err := media.DetectMIME(file)
	if err != nil {
		file.Close()
		return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unable to process file", Cause: err}
	}
	file.Seek(0, io.SeekStart)

	var kind string
	if idx := strings.Index(mime, "/"); idx != -1 {
		kind = mime[:idx]
	}

	if kind != "image" && kind != "audio" {
		file.Close()
		return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "invalid file type"}
	}

	return multipartInfo{file: file, model: model, kind: kind}, nil
}
