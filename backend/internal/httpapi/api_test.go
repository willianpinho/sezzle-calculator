package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateSuccess(t *testing.T) {
	for _, test := range []struct {
		name, body string
		want       float64
	}{
		{"add", `{"operation":"add","a":2,"b":3}`, 5},
		{"subtract", `{"operation":"subtract","a":2,"b":3}`, -1},
		{"multiply", `{"operation":"multiply","a":-2,"b":3}`, -6},
		{"divide", `{"operation":"divide","a":7,"b":2}`, 3.5},
		{"power", `{"operation":"power","a":2,"b":8}`, 256},
		{"sqrt", `{"operation":"sqrt","a":9}`, 3},
		{"percentage", `{"operation":"percentage","a":25,"b":80}`, 20},
		{"zero power zero", `{"operation":"power","a":0,"b":0}`, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := request(http.MethodPost, "/api/calculate", test.body, "application/json")
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body)
			}
			var response map[string]float64
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil || len(response) != 1 || response["result"] != test.want {
				t.Fatalf("response = %+v, err = %v", response, err)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tooLarge := `{"operation":"sqrt","a":1,"padding":"` + strings.Repeat("x", int(MaxBodyBytes)) + `"}`
	leadingTooLarge := strings.Repeat(" ", int(MaxBodyBytes)+1) + `{}`
	for _, test := range []struct {
		name, method, path, body, contentType, code string
		status                                      int
	}{
		{"malformed", "POST", "/api/calculate", `{`, "application/json", "invalid_json", 400},
		{"non object", "POST", "/api/calculate", `[]`, "application/json", "invalid_json", 400},
		{"multiple values", "POST", "/api/calculate", `{} {}`, "application/json", "invalid_json", 400},
		{"duplicate key", "POST", "/api/calculate", `{"a":1,"a":2,"operation":"sqrt"}`, "application/json", "invalid_json", 400},
		{"wrong content type", "POST", "/api/calculate", `{}`, "text/plain", "unsupported_media_type", 415},
		{"unknown field", "POST", "/api/calculate", `{"operation":"sqrt","a":1,"extra":1}`, "application/json", "validation_error", 422},
		{"case sensitive field", "POST", "/api/calculate", `{"Operation":"sqrt","a":1}`, "application/json", "validation_error", 422},
		{"case sensitive operation", "POST", "/api/calculate", `{"operation":"ADD","a":1,"b":2}`, "application/json", "validation_error", 422},
		{"missing operation", "POST", "/api/calculate", `{"a":1,"b":2}`, "application/json", "validation_error", 422},
		{"non-string operation", "POST", "/api/calculate", `{"operation":42,"a":1,"b":2}`, "application/json", "validation_error", 422},
		{"missing operand", "POST", "/api/calculate", `{"operation":"divide","a":1}`, "application/json", "validation_error", 422},
		{"null operand", "POST", "/api/calculate", `{"operation":"sqrt","a":null}`, "application/json", "validation_error", 422},
		{"null operation", "POST", "/api/calculate", `{"operation":null,"a":1}`, "application/json", "validation_error", 422},
		{"null binary operand", "POST", "/api/calculate", `{"operation":"add","a":1,"b":null}`, "application/json", "validation_error", 422},
		{"string operand", "POST", "/api/calculate", `{"operation":"sqrt","a":"1"}`, "application/json", "validation_error", 422},
		{"quoted numeric operand", "POST", "/api/calculate", `{"operation":"sqrt","a":"2"}`, "application/json", "validation_error", 422},
		{"sqrt b", "POST", "/api/calculate", `{"operation":"sqrt","a":1,"b":0}`, "application/json", "validation_error", 422},
		{"number overflow", "POST", "/api/calculate", `{"operation":"sqrt","a":1e400}`, "application/json", "validation_error", 422},
		{"divide zero", "POST", "/api/calculate", `{"operation":"divide","a":1,"b":0}`, "application/json", "domain_error", 422},
		{"negative sqrt", "POST", "/api/calculate", `{"operation":"sqrt","a":-1}`, "application/json", "domain_error", 422},
		{"fractional negative power", "POST", "/api/calculate", `{"operation":"power","a":-1,"b":0.5}`, "application/json", "domain_error", 422},
		{"result overflow", "POST", "/api/calculate", `{"operation":"multiply","a":1.7976931348623157e308,"b":2}`, "application/json", "domain_error", 422},
		{"payload too large", "POST", "/api/calculate", tooLarge, "application/json", "payload_too_large", 413},
		{"leading payload too large", "POST", "/api/calculate", leadingTooLarge, "application/json", "payload_too_large", 413},
		{"wrong method", "GET", "/api/calculate", "", "", "method_not_allowed", 405},
		{"unknown route", "GET", "/api/nope", "", "", "not_found", 404},
		{"api namespace root", "GET", "/api", "", "", "not_found", 404},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := request(test.method, test.path, test.body, test.contentType)
			if recorder.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.status, recorder.Body)
			}
			if got := errorCode(t, recorder); got != test.code {
				t.Fatalf("error code = %q, want %q", got, test.code)
			}
			if contentType := recorder.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
				t.Fatalf("content type = %q", contentType)
			}
		})
	}
}

func TestHealthAndStatic(t *testing.T) {
	static := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("static path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte("ui"))
	})
	handler := New(static)
	for _, test := range []struct {
		path string
		want int
	}{{"/api/health", 200}, {"/", 200}} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, test.path, nil))
		if recorder.Code != test.want {
			t.Errorf("GET %s: %d", test.path, recorder.Code)
		}
	}
}

func request(method, path, body, contentType string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	recorder := httptest.NewRecorder()
	New(nil).ServeHTTP(recorder, req)
	return recorder
}

func errorCode(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()
	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response.Error.Code
}
