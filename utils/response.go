package utils

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/satoshi03/related-article-api/common"
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
	w.Write([]byte(fmt.Sprintf("%s(%s)", common.JsonpNameSpace, jsonBuffer.String())))
}
