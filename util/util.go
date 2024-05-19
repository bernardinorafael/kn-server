package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseBodyJSON(r *http.Response, payload any) error {
	if r.Body == nil {
		return errors.New("missing body request")
	}

	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		return err
	}

	return nil
}
