package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthMethodAndNoStaticFallback(t *testing.T) {
	for _, test := range []struct {
		handler http.Handler
		path    string
		want    int
		code    string
	}{
		{New(nil), "/api/health", http.StatusMethodNotAllowed, "method_not_allowed"},
		{New(nil), "/missing", http.StatusNotFound, "not_found"},
	} {
		recorder := httptest.NewRecorder()
		test.handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, test.path, nil))
		if recorder.Code != test.want || errorCode(t, recorder) != test.code {
			t.Fatalf("%s: status=%d body=%s", test.path, recorder.Code, recorder.Body)
		}
	}
}

func TestDecoderBoundaryErrors(t *testing.T) {
	for _, body := range []string{`{"operation":"sqrt","a":1`, `{"a":1]`, `{} !`} {
		recorder := request(http.MethodPost, "/api/calculate", body, "application/json")
		if recorder.Code != http.StatusBadRequest || errorCode(t, recorder) != "invalid_json" {
			t.Fatalf("body=%q status=%d response=%s", body, recorder.Code, recorder.Body)
		}
	}
}
