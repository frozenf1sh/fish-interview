package web

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frozenf1sh/fish-interview/internal/content"
)

func newDemoServer(t *testing.T) *Server {
	t.Helper()
	catalog, err := content.Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatal(err)
	}
	server, err := New(catalog)
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestHomeSearchesCards(t *testing.T) {
	server := newDemoServer(t)
	req := httptest.NewRequest(http.MethodGet, "/?q=Rebalance", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	if !strings.Contains(res.Body.String(), "Kafka Rebalance") {
		t.Fatal("home should include matching card")
	}
}

func TestCardRendersInternalLink(t *testing.T) {
	server := newDemoServer(t)
	req := httptest.NewRequest(http.MethodGet, "/cards/eng.kafka.consumer-group", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	if !strings.Contains(res.Body.String(), `href="/cards/eng.kafka.rebalance"`) {
		t.Fatal("card should render wikilink")
	}
}

func TestTraceEndpointReturnsReplayableFrames(t *testing.T) {
	server := newDemoServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/traces/interval-scheduling", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"frames"`) {
		t.Fatalf("trace response = %d %s", res.Code, res.Body.String())
	}
}
