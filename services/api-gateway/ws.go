package main

import (
	"log"
	"net/http"

	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("no userID is provided")
		return

	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading the message: %v", err)
			break
		}

		log.Printf("Received message %v ", []byte(message))
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("no userID is provided")
		return

	}
	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Println("no packageSlug is provided")
		return

	}

	type Driver struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:             userID,
			Name:           "adi",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "abc123",
			PackageSlug:    packageSlug,
		},
	}
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("error sending message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading the message: %v", err)
			break
		}

		log.Printf("Received message %v ", []byte(message))
	}
}
