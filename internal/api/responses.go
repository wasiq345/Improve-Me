package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) error {
	resp, err := json.Marshal(payload)

	if err != nil {
		return errors.New("Unable to Marshal")
	}

	w.WriteHeader(code)
	w.Write(resp)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, customError string) error {
	err := RespondWithJson(w, code, map[string]string{"error: ": customError})
	if err != nil {
		return errors.New("Error in RespondWithError")
	}
	return nil
}
