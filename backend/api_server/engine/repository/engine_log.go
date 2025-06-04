package repository

import "time"

type EngineLog struct {
	ID         int       `json:"id"`
	ModelingID int       `json:"modeling_id"`
	Level      string    `json:"level"`
	Filename   string    `json:"filename"`
	Line       int       `json:"line"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
