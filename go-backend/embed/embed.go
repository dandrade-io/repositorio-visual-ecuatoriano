package embedservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repository/config"
	"repository/database"
)

func Text(ctx context.Context, text string, model database.Model) (embedTextResp, error) {
	body := &bytes.Buffer{}
	payload := &embedTextReq{Text: text, Model: model.String()}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return embedTextResp{}, err
	}

	url := config.EmbedServiceURL + "/api/v1/embed/text"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return embedTextResp{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return embedTextResp{}, fmt.Errorf("failed to make http petition: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return embedTextResp{}, fmt.Errorf("failed to make request (%d): %s", resp.StatusCode, string(b))
	}

	var embedResp embedTextResp
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return embedTextResp{}, err
	}

	return embedResp, nil

}
