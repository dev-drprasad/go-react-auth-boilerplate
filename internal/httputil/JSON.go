package httputil

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Message string `json:"message"`
}

const (
	charsetUTF8                    = "charset=UTF-8"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
)

func JSON(w http.ResponseWriter, r *http.Request, code int, i interface{}) (err error) {
	enc := json.NewEncoder(w)
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	return enc.Encode(i)
}

func JSONError(w http.ResponseWriter, r *http.Request, code int, err error) error {
	return JSONMessage(w, r, code, err.Error())
}

func JSONMessage(w http.ResponseWriter, r *http.Request, code int, message string) error {
	return JSON(w, r, code, jsonError{Message: message})
}
