package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var t Template
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		event, err := Publish(t)
		status := http.StatusOK
		if err != nil {
			status = http.StatusUnprocessableEntity
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(event)
	})
	log.Println("sms template service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
