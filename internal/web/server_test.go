package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frozenf1sh/fish-interview/internal/content"
	"github.com/frozenf1sh/fish-interview/internal/trace"
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

func TestHealthEndpoints(t *testing.T) {
	server := newDemoServer(t)
	for _, path := range []string{"/healthz", "/readyz"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			res := httptest.NewRecorder()
			server.ServeHTTP(res, req)
			if res.Code != http.StatusOK || res.Body.String() != "ok\n" {
				t.Fatalf("health response for %s = %d %q", path, res.Code, res.Body.String())
			}
		})
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

func TestPatternTraceEndpoints(t *testing.T) {
	server := newDemoServer(t)
	for _, path := range []string{
		"/api/traces/linear-dp", "/api/traces/space-rolling", "/api/traces/lcs-dp", "/api/traces/interval-dp",
		"/api/traces/stock-dp", "/api/traces/bitmask-dp", "/api/traces/linked-list-rewire", "/api/traces/path-dp", "/api/traces/reverse-path-dp", "/api/traces/binary-red-blue",
		"/api/traces/flow-bfs-shortest-path", "/api/traces/flow-tree-dp", "/api/traces/flow-string-golang",
		"/api/traces/lis", "/api/traces/row-gravity",
		"/api/traces/interval-start-merge", "/api/traces/meeting-rooms", "/api/traces/weighted-intervals", "/api/traces/kadane",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"frames"`) {
			t.Fatalf("trace response for %s = %d %s", path, res.Code, res.Body.String())
		}
	}
}

func TestEveryTraceEndpointMeetsPlayerContract(t *testing.T) {
	server := newDemoServer(t)
	for name := range trace.AllTraces() {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/traces/"+name, nil)
			res := httptest.NewRecorder()
			server.ServeHTTP(res, req)
			if res.Code != http.StatusOK {
				t.Fatalf("trace response = %d: %s", res.Code, res.Body.String())
			}

			var playerTrace trace.Trace
			if err := json.Unmarshal(res.Body.Bytes(), &playerTrace); err != nil {
				t.Fatalf("decode trace response: %v", err)
			}
			if err := trace.ValidatePlayerContract(playerTrace); err != nil {
				t.Fatalf("invalid trace response: %v", err)
			}
		})
	}
}

func TestPatternCardsEmbedKnownTraces(t *testing.T) {
	server := newDemoServer(t)
	traces := map[string]string{
		"algo.greedy.interval-scheduling":  "/api/traces/interval-scheduling",
		"algo.dp.linear":                   "/api/traces/linear-dp",
		"algo.dp.lcs":                      "/api/traces/lcs-dp",
		"algo.dp.interval":                 "/api/traces/interval-dp",
		"algo.dp.stock":                    "/api/traces/stock-dp",
		"algo.dp.bitmask":                  "/api/traces/bitmask-dp",
		"algo.list.dummy-rewire":           "/api/traces/linked-list-rewire",
		"algo.dp.path":                     "/api/traces/path-dp",
		"algo.dp.path.minimum-health":      "/api/traces/reverse-path-dp",
		"algo.binary-search.answer":        "/api/traces/binary-red-blue",
		"algo.bfs.shortest-path":           "/api/traces/flow-bfs-shortest-path",
		"algo.tree.dp":                     "/api/traces/flow-tree-dp",
		"algo.string.golang":               "/api/traces/flow-string-golang",
		"algo.sequence.lis":                "/api/traces/lis",
		"algo.simulation.gravity":          "/api/traces/row-gravity",
		"algo.greedy.interval-start-merge": "/api/traces/interval-start-merge",
		"algo.greedy.meeting-rooms":        "/api/traces/meeting-rooms",
		"algo.greedy.weighted-intervals":   "/api/traces/weighted-intervals",
		"algo.greedy.kadane":               "/api/traces/kadane",
	}
	for id, traceURL := range traces {
		t.Run(id, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cards/"+id, nil)
			res := httptest.NewRecorder()
			server.ServeHTTP(res, req)
			if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `data-trace="`+traceURL+`"`) {
				t.Fatalf("card %s did not embed %s: %d", id, traceURL, res.Code)
			}
		})
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

func TestCardIncludesGestureTreeControls(t *testing.T) {
	server := newDemoServer(t)
	req := httptest.NewRequest(http.MethodGet, "/cards/algo.dp.linear?tree=algorithms", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	body := res.Body.String()
	if res.Code != http.StatusOK || !strings.Contains(body, "双指缩放") || !strings.Contains(body, "data-tree-focus") || strings.Contains(body, "data-tree-zoom-in") {
		t.Fatalf("card should include tree controls: %d", res.Code)
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

func TestMarkdownRendersLatexStateTransition(t *testing.T) {
	server := newDemoServer(t)
	rendered, err := server.renderMarkdown("$$dp_{i} = dp_{i-1} + dp_{i-2}$$")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rendered, `class="math-formula"`) || !strings.Contains(rendered, `dp<sub>i</sub>`) {
		t.Fatalf("expected rendered formula, got %s", rendered)
	}
}

func TestDPCardsRenderFormulaAndSegmentedImplementation(t *testing.T) {
	server := newDemoServer(t)
	for _, id := range []string{"algo.dp.linear", "algo.dp.lcs", "algo.dp.interval", "algo.dp.stock", "algo.dp.bitmask", "algo.dp.path", "algo.dp.path.minimum-health"} {
		t.Run(id, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cards/"+id, nil)
			res := httptest.NewRecorder()
			server.ServeHTTP(res, req)
			body := res.Body.String()
			if res.Code != http.StatusOK || !strings.Contains(body, `class="math-formula"`) || !strings.Contains(body, "分段实现") {
				t.Fatalf("card %s misses formula or segmented implementation: %d", id, res.Code)
			}
		})
	}
}
