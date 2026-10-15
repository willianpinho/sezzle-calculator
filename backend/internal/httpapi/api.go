// Package httpapi exposes the calculator's HTTP contract.
package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"strconv"

	"github.com/willianpinho/sezzle-calculator/backend/internal/calculator"
)

const MaxBodyBytes int64 = 1 << 20

type API struct{ static http.Handler }

func New(static http.Handler) http.Handler { return API{static: static} }

func (a API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/calculate":
		a.calculate(w, r)
	case "/api/health":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	default:
		if r.URL.Path == "/api" || (len(r.URL.Path) >= 5 && r.URL.Path[:5] == "/api/") {
			writeError(w, http.StatusNotFound, "not_found", "route not found")
			return
		}
		if a.static != nil {
			a.static.ServeHTTP(w, r)
			return
		}
		writeError(w, http.StatusNotFound, "not_found", "route not found")
	}
}

func (a API) calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "content type must be application/json")
		return
	}
	fields, err := decodeObject(w, r)
	if err != nil {
		return
	}
	op, err := requiredString(fields, "operation")
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}
	if !validOperation(op) {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", "operation is invalid")
		return
	}
	aValue, err := requiredNumber(fields, "a")
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}
	var bValue float64
	if op == "sqrt" {
		if _, ok := fields["b"]; ok {
			writeError(w, http.StatusUnprocessableEntity, "validation_error", "b is not allowed for sqrt")
			return
		}
	} else if bValue, err = requiredNumber(fields, "b"); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}
	result, err := calculator.Calculate(op, aValue, bValue)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "domain_error", "calculation has no finite real result")
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"result": result})
}

func decodeObject(w http.ResponseWriter, r *http.Request) (map[string]json.RawMessage, error) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	decoder := json.NewDecoder(r.Body)
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		if err != nil {
			writeDecodeError(w, err)
			return nil, err
		}
		writeError(w, http.StatusBadRequest, "invalid_json", "body must be one JSON object")
		return nil, errors.New("body is not an object")
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		keyToken, decodeErr := decoder.Token()
		// JSON grammar tokens/errors are public request syntax, never credentials.
		if decodeErr != nil {
			writeDecodeError(w, decodeErr)
			return nil, decodeErr
		}
		key, ok := keyToken.(string)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid_json", "malformed JSON")
			return nil, errors.New("invalid key")
		}
		if _, exists := fields[key]; exists {
			writeError(w, http.StatusBadRequest, "invalid_json", "duplicate object field")
			return nil, errors.New("duplicate field")
		}
		var value json.RawMessage
		if decodeErr = decoder.Decode(&value); decodeErr != nil {
			writeDecodeError(w, decodeErr)
			return nil, decodeErr
		}
		fields[key] = value
	}
	if _, err = decoder.Token(); err != nil {
		writeDecodeError(w, err)
		return nil, err
	}
	if err = ensureEOF(decoder); err != nil {
		writeDecodeError(w, err)
		return nil, err
	}
	for key := range fields {
		if key != "operation" && key != "a" && key != "b" {
			writeError(w, http.StatusUnprocessableEntity, "validation_error", "unknown field: "+key)
			return nil, errors.New("unknown field")
		}
	}
	return fields, nil
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	err := decoder.Decode(&extra)
	if err == io.EOF {
		return nil
	}
	if err == nil {
		return errors.New("multiple JSON values")
	}
	return err
}
func requiredString(fields map[string]json.RawMessage, key string) (string, error) {
	raw, ok := fields[key]
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	var value string
	if json.Unmarshal(raw, &value) != nil {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return value, nil
}
func requiredNumber(fields map[string]json.RawMessage, key string) (float64, error) {
	raw, ok := fields[key]
	if !ok {
		return 0, fmt.Errorf("%s is required", key)
	}
	if len(raw) == 0 || (raw[0] != '-' && (raw[0] < '0' || raw[0] > '9')) {
		return 0, fmt.Errorf("%s must be a number", key)
	}
	var number json.Number
	if json.Unmarshal(raw, &number) != nil {
		return 0, fmt.Errorf("%s must be a number", key)
	}
	value, err := strconv.ParseFloat(number.String(), 64)
	if err != nil || math.IsInf(value, 0) || math.IsNaN(value) {
		return 0, fmt.Errorf("%s must be a finite number", key)
	}
	return value, nil
}
func validOperation(op string) bool {
	switch op {
	case "add", "subtract", "multiply", "divide", "power", "sqrt", "percentage":
		return true
	}
	return false
}
func writeDecodeError(w http.ResponseWriter, err error) {
	var maxErr *http.MaxBytesError
	if errors.As(err, &maxErr) {
		writeError(w, http.StatusRequestEntityTooLarge, "payload_too_large", "request body is too large")
		return
	}
	writeError(w, http.StatusBadRequest, "invalid_json", "malformed JSON")
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
