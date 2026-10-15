package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestOpenAPIContract(t *testing.T) {
	recorder := request(http.MethodGet, "/api/openapi.json", "", "")
	if recorder.Code != http.StatusOK || !bytes.Equal(recorder.Body.Bytes(), openAPIContract) {
		t.Fatalf("status=%d body=%q", recorder.Code, recorder.Body)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Fatalf("content type = %q", contentType)
	}
	var contract map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &contract); err != nil {
		t.Fatal(err)
	}
	if contract["openapi"] != "3.1.1" {
		t.Fatalf("openapi = %#v", contract["openapi"])
	}
	paths, ok := contract["paths"].(map[string]any)
	if !ok {
		t.Fatal("paths is missing")
	}
	for _, path := range []string{"/api/calculate", "/api/health", "/api/openapi.json", "/api/docs"} {
		if _, ok := paths[path]; !ok {
			t.Errorf("missing path %s", path)
		}
	}
	validateReferences(t, contract, contract)
}

func TestDocumentationRoutesRejectOtherMethods(t *testing.T) {
	for _, path := range []string{"/api/openapi.json", "/api/docs"} {
		recorder := request(http.MethodPost, path, "", "")
		if recorder.Code != http.StatusMethodNotAllowed || recorder.Header().Get("Allow") != http.MethodGet || errorCode(t, recorder) != "method_not_allowed" {
			t.Fatalf("POST %s: status=%d body=%s", path, recorder.Code, recorder.Body)
		}
	}
}

func TestDocsRedirect(t *testing.T) {
	recorder := request(http.MethodGet, "/api/docs", "", "")
	if recorder.Code != http.StatusTemporaryRedirect || recorder.Header().Get("Location") != "/docs/" {
		t.Fatalf("status=%d location=%q", recorder.Code, recorder.Header().Get("Location"))
	}
}

func TestDocumentationNamespaceIsJSONNotFound(t *testing.T) {
	for _, path := range []string{"/api/openapi.json/extra", "/api/docs/extra"} {
		recorder := request(http.MethodGet, path, "", "")
		if recorder.Code != http.StatusNotFound || errorCode(t, recorder) != "not_found" {
			t.Fatalf("GET %s: status=%d body=%s", path, recorder.Code, recorder.Body)
		}
	}
}

func validateReferences(t *testing.T, node any, root map[string]any) {
	t.Helper()
	switch value := node.(type) {
	case map[string]any:
		if reference, ok := value["$ref"].(string); ok {
			resolveReference(t, reference, root)
		}
		for _, child := range value {
			validateReferences(t, child, root)
		}
	case []any:
		for _, child := range value {
			validateReferences(t, child, root)
		}
	}
}

func resolveReference(t *testing.T, reference string, root map[string]any) {
	t.Helper()
	if !strings.HasPrefix(reference, "#/") {
		t.Errorf("unsupported reference %q", reference)
		return
	}
	var value any = root
	for _, segment := range strings.Split(strings.TrimPrefix(reference, "#/"), "/") {
		object, ok := value.(map[string]any)
		if !ok {
			t.Errorf("reference %q does not resolve", reference)
			return
		}
		value, ok = object[strings.ReplaceAll(strings.ReplaceAll(segment, "~1", "/"), "~0", "~")]
		if !ok {
			t.Errorf("reference %q does not resolve", reference)
			return
		}
	}
}
