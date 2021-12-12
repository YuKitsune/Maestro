package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResource struct {
	Error string
}

func Response(w http.ResponseWriter, res interface{}, status int) {

	resBytes, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		Error(w, err)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(resBytes)
	if err != nil {
		Error(w, err)
		return
	}
}

func NotFound(w http.ResponseWriter, message string) {
	res := &ErrorResource{message}
	Response(w, res, http.StatusNotFound)
}

func NotFoundf(w http.ResponseWriter, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v)
	res := &ErrorResource{msg}
	Response(w, res, http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, message string) {
	res := &ErrorResource{message}
	Response(w, res, http.StatusBadRequest)
}

func Error(w http.ResponseWriter, err error) {
	res := &ErrorResource{err.Error()}
	Response(w, res, http.StatusInternalServerError)
}
