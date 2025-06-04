package dto

type WebSocketDTO struct {
    MessageType string                  `json:"messageType"`
    Data        map[string]interface{}  `json:"data"`
}
