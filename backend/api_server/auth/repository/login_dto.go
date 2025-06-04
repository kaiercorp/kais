package repository

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Username string    `json:"username,omitempty"`
	Token    string    `json:"token,omitempty"`
	LoginAt  time.Time `json:"login_at,omitempty"`
	Name     string    `json:"name,omitempty"`
	Group    *int      `json:"group,omitempty"`
}

type LogoutRequest struct {
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}
