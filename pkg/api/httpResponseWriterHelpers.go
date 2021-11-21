package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Response(w http.ResponseWriter, res interface{}, status int) {

	resBytes, err := json.Marshal(res)
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

func EmptyResponse(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func BadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, message)
}

func BadRequestf(w http.ResponseWriter, format string, s ...interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, format, s...)
}

func Error(w http.ResponseWriter, err error) {
	switch err.(type) {

	// Todo: Write known errors here

	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		break
	}
}

func EmptyError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
