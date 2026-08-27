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

func TestCardPartialEmbedsAlgorithmTraceWithoutLayout(t *testing.T) {
	server := newDemoServer(t)
	req := httptest.NewRequest(http.MethodGet, "/partials/cards/algo.greedy.interval-scheduling?tree=algorithms", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	body := res.Body.String()
	if !strings.Contains(body, `data-trace="/api/traces/interval-scheduling"`) {
		t.Fatal("partial should embed the algorithm trace player")
	}
	if strings.Contains(body, "<html") || strings.Contains(body, "tree-canvas") {
		t.Fatal("partial should replace only the content panel")
	}
}

func TestMarkdownUsesChromaWithLineNumbers(t *testing.T) {
	server := newDemoServer(t)
	rendered, err := server.renderMarkdown("```go\nfmt.Println(\"fish\")\n```")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered, `class="code-block"`) || !strings.Contains(rendered, "user-select:none") || !strings.Contains(rendered, "color:#") {
		t.Fatalf("expected highlighted line-numbered output, got %s", rendered)
	}
}

func TestMarkdownKeepsChineseEmphasisBeforeColon(t *testing.T) {
	server := newDemoServer(t)
	rendered, err := server.renderMarkdown("**重点：**后面是说明")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered, "<strong>重点</strong>：后面是说明") {
		t.Fatalf("expected normalized emphasis, got %s", rendered)
	}
}
