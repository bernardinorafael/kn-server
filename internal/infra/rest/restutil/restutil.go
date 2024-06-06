package restutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteSuccess(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func ParseBody(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(v)
}
