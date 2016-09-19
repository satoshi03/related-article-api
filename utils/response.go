package utils

import (
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
)

func Write404Response(w http.ResponseWriter, msg map[string]interface{}) {
	WriteResponse(w, msg, 404)
}

func WriteResponse(w http.ResponseWriter, msg map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	ffjson.NewEncoder(w).Encode(msg)
}
