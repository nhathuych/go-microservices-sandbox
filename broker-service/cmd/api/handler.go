package main

import (
	"net/http"
)

// Deep Health Check: kiểm tra coi App Go này có kết nối được với thế giới bên ngoài (RabbitMQ) không
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
