package middleware

import (
	"encoding/json"
	"errors"
	"http-proxy/internal/models"
	"net/http"
)

var (
	BadRequest          = errors.New("bad request")
	NotFound            = errors.New("item is not found")
	Conflict            = errors.New("already exist")
	InternalServerError = errors.New("internal Server Error")
)

func JsonError(message string) []byte {
	jsonErr, err := json.Marshal(models.Error{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonErr
}

func Response(w http.ResponseWriter, status models.StatusCode, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	switch status {
	case models.Okey:
		w.WriteHeader(http.StatusOK)
	case models.Created:
		w.WriteHeader(http.StatusCreated)
	case models.NotFound:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(JsonError(NotFound.Error()))
		return
	case models.Conflict:
		w.WriteHeader(http.StatusConflict)
		if body != nil {
			break
		} else {
			_, _ = w.Write(JsonError(Conflict.Error()))
		}
		return
	case models.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(JsonError(BadRequest.Error()))
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(JsonError(InternalServerError.Error()))
		return
	}

	if body != nil {
		jsn, err := json.Marshal(body)
		if err != nil {
			return
		}
		_, _ = w.Write(jsn)
	}
}
