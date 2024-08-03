package routeutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteSuccessResponse(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func WriteJSONResponse(w http.ResponseWriter, code int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func ParseBodyRequest(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("invalid request data")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
