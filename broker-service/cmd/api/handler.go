package main

import (
	"encoding/json"
	"net/http"
)

// Deep Health Check: kiểm tra coi App Go này có kết nối được với thế giới bên ngoài (RabbitMQ) không
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
