package main

import (
	"encoding/json"
	"net/http"

	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	// call trip service
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse the request body", http.StatusBadRequest)
		return
	}
	if reqBody.UserID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	response := contracts.APIResponse{Data: "ok"}
	writeJSON(w, http.StatusCreated, response)
}
