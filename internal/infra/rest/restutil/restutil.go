package restutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteSuccess(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func ParseBody(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("invalid request data")
	}
	return json.NewDecoder(r.Body).Decode(v)
}
