package repository

import (
	"database/sql"
	"time"
)

type EngineLogResponse struct {
	ID    int            `json:"id"`
	Level sql.NullString `json:"level"`
	// FileName  sql.NullString `json:"filename"
	Line         sql.NullString `json:"line"`
	Message      string         `json:"message"`
	CreatedAt    time.Time      `json:"created_at"`
	ModelingID   int            `json:"modeling_id"`
	ModelingStep string         `json:"modeling_step"`
}
