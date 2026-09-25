package assets

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"repository/database"
	"repository/files"
	transporthttp "repository/transport/http"
	"strconv"
	"strings"
	"uuid"
)

type Handler struct {
	Service Service
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) error {
	embed := true
	if val, err := strconv.ParseBool(r.URL.Query().Get("embed")); err != nil {
		embed = val
	}

	info, err := parseMultipartFile(r, embed)
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		return err
	}
	defer info.file.Close()

	asset, err := h.Service.Create(info.file, info.kind)
	if err != nil {
		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to store asset", Cause: err}
	}

	if embed {
		if err := h.Service.embed(r.Context(), asset, info.model); err != nil {
			return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to embed asset", Cause: err}
		}
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

	model, err := database.ParseModel(req.Model)
	if err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model"}
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "invalid asset id"}
	}

	if err := h.Service.EmbedByAssetID(r.Context(), id, model); err != nil {
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

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("q")
	kind := r.URL.Query().Get("kind")
	model := database.JinaClipV2
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
		val, err := database.ParseModel(name)
		if err != nil {
			return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model provided"}
		}
		model = val
	}
	if num, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = num
	}

	embeddings, err := h.Service.search(r.Context(), q, model, limit)
	if err != nil {
		return &transporthttp.Error{Status: http.StatusInternalServerError, Msg: "unable to search for k-best coincidences", Cause: err}
	}

	result := searchResp{Results: make([]searchResult, len(embeddings)), NumResults: len(embeddings)}
	for i := range embeddings {
		e := &embeddings[i]
		result.Results[i] = searchResult{AssetId: e.AssetId, Similarity: e.Similarity}
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
	model database.Model
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

	var model database.Model
	if embed {
		if len(models) != 1 {
			return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "only one model value must be provided"}
		}
		val, err := database.ParseModel(models[0])
		if err != nil {
			return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unknown model provided"}
		}
		model = val
	}

	header := headers[0]
	file, err := header.Open()
	if err != nil {
		return multipartInfo{}, &transporthttp.Error{Status: http.StatusBadRequest, Msg: "unable to process file", Cause: err}
	}

	mime, err := mimeType(file)
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
