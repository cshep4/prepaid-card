package response

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, code int, msg string) {
	Json(w, code, map[string]string{"error": msg})
}

func Json(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

