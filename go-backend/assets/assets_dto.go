package assets

import "uuid"

type assetResp struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"path"`
}

type embedReq struct {
	Model string `json:"model"`
}
