package httpapi

import (
	"encoding/json"
	"math"
	"net/http"
	"testing"
)

func TestHTTPNumericBoundaries(t *testing.T) {
	for _, test := range []struct {
		name, body, code string
		status           int
	}{
		{"signed zero division", `{"operation":"divide","a":1,"b":-0}`, "domain_error", 422},
		{"zero over zero", `{"operation":"divide","a":0,"b":0}`, "domain_error", 422},
		{"zero negative power", `{"operation":"power","a":0,"b":-2}`, "domain_error", 422},
		{"negative fractional power", `{"operation":"power","a":-8,"b":0.3333333333333333}`, "domain_error", 422},
		{"add overflow", `{"operation":"add","a":1.7976931348623157e308,"b":1.7976931348623157e308}`, "domain_error", 422},
		{"subtract overflow", `{"operation":"subtract","a":1.7976931348623157e308,"b":-1.7976931348623157e308}`, "domain_error", 422},
		{"multiply overflow", `{"operation":"multiply","a":1.7976931348623157e308,"b":2}`, "domain_error", 422},
		{"divide overflow", `{"operation":"divide","a":1.7976931348623157e308,"b":5e-324}`, "domain_error", 422},
		{"power overflow", `{"operation":"power","a":1.7976931348623157e308,"b":2}`, "domain_error", 422},
		{"percentage overflow", `{"operation":"percentage","a":1.7976931348623157e308,"b":1.7976931348623157e308}`, "domain_error", 422},
		{"positive overflow input", `{"operation":"sqrt","a":1e400}`, "validation_error", 422},
		{"negative overflow input", `{"operation":"sqrt","a":-1e400}`, "validation_error", 422},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := request(http.MethodPost, "/api/calculate", test.body, "application/json")
			if recorder.Code != test.status || errorCode(t, recorder) != test.code {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body)
			}
		})
	}
}

func TestHTTPFiniteNumericSuccess(t *testing.T) {
	for _, test := range []struct {
		name, body string
		want       float64
		delta      float64
	}{
		{"negative odd integer power", `{"operation":"power","a":-2,"b":3}`, -8, 0},
		{"positive fractional power", `{"operation":"power","a":9,"b":0.5}`, 3, 0},
		{"subnormal square root", `{"operation":"sqrt","a":5e-324}`, math.Sqrt(math.SmallestNonzeroFloat64), 0},
		{"rounding at two to fifty-three", `{"operation":"add","a":9007199254740992,"b":1}`, 1 << 53, 0},
		{"percentage underflow", `{"operation":"percentage","a":5e-324,"b":50}`, 0, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := request(http.MethodPost, "/api/calculate", test.body, "application/json")
			if recorder.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body)
			}
			var response map[string]float64
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil || len(response) != 1 || math.Abs(response["result"]-test.want) > test.delta {
				t.Fatalf("response=%v err=%v", response, err)
			}
		})
	}
}
