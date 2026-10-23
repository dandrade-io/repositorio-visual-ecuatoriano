package albumcover

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"repository/database"
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

type uploadReq struct {
	Title         string    `json:"title"`
	Year          *int      `json:"year"`
	Artist        string    `json:"artist"`
	Rhythms       []string  `json:"rhythms"`
	ReprTechnique string    `json:"representationTechnique"`
	Predominance  string    `json:"graphicPredominance"`
	Colors        []string  `json:"dominantColors"`
	TitleSyntax   string    `json:"titleSyntax"`
	AssetID       uuid.UUID `json:"assetId"`
}

func (r uploadReq) toInput(userID uuid.UUID) CreateInput {
	return CreateInput{
		Title:         r.Title,
		Year:          r.Year,
		Artist:        r.Artist,
		Rhythms:       r.Rhythms,
		ReprTechnique: r.ReprTechnique,
		Predominance:  r.Predominance,
		Colors:        r.Colors,
		TitleSyntax:   r.TitleSyntax,
		AssetID:       r.AssetID,
		PublishedBy:   userID,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	type resp struct {
		Results []database.AlbumCoverSummary `json:"results"`
	}

	limit := 10
	if num, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = num
	}

	covers, err := h.Svc.FetchSummary(r.Context(), limit)
	if err != nil {
		return err
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&resp{Results: covers}); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) error {
	claims, _ := transporthttp.ClaimsFrom(r.Context())
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return err
	}

	var req uploadReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	ac, err := h.Svc.Create(r.Context(), req.toInput(userID))
	if err != nil {
		return err
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&ac); err != nil {
		return err
	}

	return nil
}

func (h *Handler) MetadataFromUpload(w http.ResponseWriter, r *http.Request) error {
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

func (h *Handler) MetadataFromAsset(w http.ResponseWriter, r *http.Request) error {
	assetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return &transporthttp.Error{Status: http.StatusBadRequest, Msg: "invalid asset id"}
	}

	file, err := os.Open(files.GetAbsPath(filepath.Join("image", assetID.String())))
	if err != nil {
		return err
	}
	defer file.Close()

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
