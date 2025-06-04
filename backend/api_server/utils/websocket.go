package utils

import (
	"net/http"
	"github.com/gorilla/websocket"
)

var WSUpgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

