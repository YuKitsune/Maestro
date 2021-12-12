package handlers

import (
	"encoding/json"
	"maestro/pkg/api/errors"
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
	Response(w, res, http.StatusBadRequest)
}

func BadRequest(w http.ResponseWriter, message string) {
	res := &ErrorResource{message}
	Response(w, res, http.StatusBadRequest)
}

func Error(w http.ResponseWriter, err error) {

	switch err.(type) {

	// Todo: Write known errors here

	case *errors.NotFoundError:
		NotFound(w, err.Error())

	default:
		res := &ErrorResource{err.Error()}
		Response(w, res, http.StatusInternalServerError)
		break
	}
}
