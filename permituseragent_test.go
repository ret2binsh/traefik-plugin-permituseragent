package traefik_plugin_permituseragent_test

import (
	"context"
	permit "github.com/ret2binsh/traefik-plugin-permituseragent"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	testUserAgentRedirection(t)
}

func testUserAgentRedirection(t *testing.T) {
	cfg := &permit.Config{
		UserAgent: "testing",
		Url: "http://www.google.com",
	}

	handler, ctx := prepare(t, cfg)

	req, recorder := prepareCase(t, ctx, "http://localhost")

	// should not redirect
	req.Header.Set("User-Agent", "testing")
	handler.ServeHTTP(recorder, req)
	assertNoRedirection(t, recorder)

	// should redirect
	recorder = httptest.NewRecorder()
	req.Header.Set("User-Agent", "wrong")
	handler.ServeHTTP(recorder, req)
	assertRedirection(t, recorder, "http://www.google.com")
}

func prepare(t *testing.T, cfg *permit.Config) (http.Handler, context.Context) {
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := permit.New(ctx, next, cfg, "cond-redirect")
	if err != nil {
		t.Fatal(err)
	}

	return handler, ctx
}

func prepareCase(t *testing.T, ctx context.Context, url string) (*http.Request, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	return req, recorder
}

func assertRedirection(t *testing.T, recorder *httptest.ResponseRecorder, location string) {
	assertStatusCode(t, recorder, 302)
	assertHeader(t, recorder, "Location", location)
}

func assertNoRedirection(t *testing.T, recorder *httptest.ResponseRecorder) {
	assertStatusCode(t, recorder, 200)
	assertHeader(t, recorder, "Location", "")
}

func assertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if recorder.Code != expected {
		t.Errorf("Wrong status code. Expected: %d. Actual: %d.", expected, recorder.Code)
	}
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	actual := recorder.Header().Get(key)
	if actual != expected {
		t.Errorf("Wrong header. Expected: %s. Actual: %s", expected, actual)
	}
}
