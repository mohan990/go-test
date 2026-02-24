package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationHealthEndpoint(t *testing.T) {
	ts := httptest.NewServer(newMux())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationHelloEndpoint(t *testing.T) {
	ts := httptest.NewServer(newMux())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/hello?name=HP")
	if err != nil {
		t.Fatalf("hello request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}
