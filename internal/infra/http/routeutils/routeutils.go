package routeutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	maxRequestBodyBytes = 1_048_576 // 1mb
)

func WriteSuccessResponse(w http.ResponseWriter, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

func WriteJSONResponse(w http.ResponseWriter, code int, dst any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(dst)
}

func ReadRequestBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxRequestBodyBytes))

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			// JSON syntax error on body request
			// Offset is the exact byte where the error occurred
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			// JSON value and struct type does not match
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			// io.EOF (End of File) indicates that there are no more bytes left to read
			return errors.New("body cannot be empty")
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &invalidUnmarshalError):
			// received a non-nin pointer into Decode()
			panic(err)
		default:
			return err
		}
	}

	// calling decode again, if there's more values after JSON sent
	// this will return an io.EOF error, that means, the client sent more data
	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
