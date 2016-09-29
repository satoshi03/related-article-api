package common

import (
	"bytes"
	"fmt"
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

func WriteJsonpResponse(w http.ResponseWriter, msg map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	jsonBuffer := new(bytes.Buffer)
	ffjson.NewEncoder(jsonBuffer).Encode(msg)
	w.Write([]byte(fmt.Sprintf("%s(%s)", JsonpNameSpace, jsonBuffer.String())))
}

type ResponseWriter func(W http.ResponseWriter, resp map[string]interface{}, statusCode int)

func JsonpWriter(w http.ResponseWriter, resp map[string]interface{}, statusCode int) {
	WriteJsonpResponse(w, resp, statusCode)
}

func JsonWriter(w http.ResponseWriter, resp map[string]interface{}, statusCode int) {
	WriteResponse(w, resp, statusCode)
}
