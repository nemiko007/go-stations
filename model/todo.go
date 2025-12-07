package model

import (
	"time"
)

// Todo はTODO情報を表します。
type Todo struct {
	ID          int64     `json:"id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTODORequest は POST /todos へのリクエストです。
type CreateTODORequest struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

// CreateTODOResponse は POST /todos へのレスポンスです。
type CreateTODOResponse struct {
	TODO Todo `json:"todo"`
}

// UpdateTODORequest は PUT /todos へのリクエストです。
type UpdateTODORequest struct {
	ID          int64  `json:"id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

// ReadTODORequest は GET /todos へのリクエストです。
type ReadTODORequest struct {
	PrevID int64 `form:"prev_id"`
	Size   int64 `form:"size"`
}

// ReadTODOResponse は GET /todos へのレスポンスです。
type ReadTODOResponse struct {
	Todos []*Todo `json:"todos"`
}

// UpdateTODOResponse は PUT /todos へのレスポンスです。
type UpdateTODOResponse struct {
	TODO Todo `json:"todo"`
}

// DeleteTODORequest は DELETE /todos へのリクエストです。
type DeleteTODORequest struct {
	IDs []int64 `json:"ids"`
}

// DeleteTODOResponse は DELETE /todos へのレスポンスです。
type DeleteTODOResponse struct {
}
