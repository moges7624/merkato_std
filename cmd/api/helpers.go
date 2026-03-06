package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type envelope map[string]any

func writeJSON(
	w http.ResponseWriter,
	status int,
	data envelope,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func readJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxByteError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"request body contains badly-formed JSON (at character %d)",
				syntaxError.Offset,
			)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"request body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}
			return fmt.Errorf(
				"request body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset,
			)

		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown key %s", fieldName)

		case errors.As(err, &maxByteError):
			return fmt.Errorf(
				"request body must be not longer than %d bytes",
				maxBytes,
			)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}

	return nil
}
