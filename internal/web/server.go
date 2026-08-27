package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	stdhtml "html"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/frozenf1sh/fish-interview/internal/content"
	"github.com/frozenf1sh/fish-interview/internal/trace"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

//go:embed templates/*.html static/*
var assets embed.FS

type Server struct {
	catalog   content.Catalog
	templates *template.Template
	markdown  goldmark.Markdown
}

type homeData struct {
	navigationData
	Title      string
	Query      string
	Results    []content.Card
	Signals    []signalCount
	TotalCards int
}

type cardData struct {
	navigationData
	Title    string
	Card     content.Card
	Body     template.HTML
	Related  []content.Card
	Signals  []content.ExamSignal
	TraceURL string
}

type labData struct {
	navigationData
	Title string
}

type signalCount struct {
	Company string
	Count   int
}

type navigationData struct {
	TreeID       string
	TreeOptions  []treeOption
	TreeJSON     template.JS
	ActiveCardID string
}

type treeOption struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type treeViewNode struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	Card     string         `json:"card,omitempty"`
	Children []treeViewNode `json:"children,omitempty"`
}

func New(catalog content.Catalog) (*Server, error) {
	templates, err := template.New("pages.html").Funcs(template.FuncMap{
		"tagClass": tagClass,
		"nodeTitle": func(node content.TreeNode) string {
			if node.Title != "" {
				return node.Title
			}
			if card, ok := catalog.Find(node.Card); ok {
				return card.Title
			}
			return node.ID
		},
	}).ParseFS(assets, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.CJK,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(chromahtml.WithLineNumbers(true)),
			),
		),
	)
	return &Server{catalog: catalog, templates: templates, markdown: markdown}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/healthz":
		s.health(w)
	case r.URL.Path == "/readyz":
		s.health(w)
	case r.URL.Path == "/":
		s.home(w, r)
	case strings.HasPrefix(r.URL.Path, "/cards/"):
		s.card(w, r)
	case strings.HasPrefix(r.URL.Path, "/partials/cards/"):
		s.cardPartial(w, r)
	case r.URL.Path == "/lab/interval-scheduling":
		s.lab(w, r)
	case r.URL.Path == "/api/traces/interval-scheduling":
		s.intervalTrace(w, r)
	case r.URL.Path == "/api/traces/linear-dp":
		s.linearDPTrace(w, r)
	case r.URL.Path == "/api/traces/space-rolling":
		s.spaceRollingTrace(w, r)
	case r.URL.Path == "/api/traces/lcs-dp":
		s.lcsTrace(w, r)
	case r.URL.Path == "/api/traces/interval-dp":
		s.intervalDPTrace(w, r)
	case r.URL.Path == "/api/traces/stock-dp":
		s.stockTrace(w, r)
	case r.URL.Path == "/api/traces/bitmask-dp":
		s.bitmaskTrace(w, r)
	case r.URL.Path == "/api/traces/linked-list-rewire":
		s.linkedListRewireTrace(w, r)
	case r.URL.Path == "/api/traces/path-dp":
		s.pathTrace(w, r)
	case r.URL.Path == "/api/traces/reverse-path-dp":
		s.reversePathTrace(w, r)
	case r.URL.Path == "/api/traces/binary-red-blue":
		s.binaryRedBlueTrace(w, r)
	case r.URL.Path == "/api/traces/lis":
		s.lisTrace(w, r)
	case r.URL.Path == "/api/traces/row-gravity":
		s.rowGravityTrace(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/traces/flow-"):
		s.flowTrace(w, r)
	case r.URL.Path == "/static/app.css":
		s.static(w, r, "static/app.css", "text/css; charset=utf-8")
	case r.URL.Path == "/static/app.js":
		s.static(w, r, "static/app.js", "text/javascript; charset=utf-8")
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) health(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}

func (s *Server) lab(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/cards/algo.greedy.interval-scheduling?tree=algorithms", http.StatusFound)
}

func (s *Server) intervalTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.IntervalScheduling([]trace.Interval{
		{Label: "A", Start: 1, End: 3},
		{Label: "B", Start: 2, End: 5},
		{Label: "C", Start: 3, End: 6},
		{Label: "D", Start: 5, End: 7},
		{Label: "E", Start: 6, End: 8},
	}))
}

func (s *Server) linearDPTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.LinearDPClimbStairs())
}

func (s *Server) spaceRollingTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.SpaceOptimizationTrace())
}

func (s *Server) lcsTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.LCSTrace())
}

func (s *Server) intervalDPTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.IntervalMergeTrace())
}

func (s *Server) stockTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.StockTrace())
}

func (s *Server) bitmaskTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.BitmaskTrace())
}

func (s *Server) linkedListRewireTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.LinkedListRewireTrace())
}

func (s *Server) pathTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.PathTrace())
}

func (s *Server) reversePathTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.ReversePathTrace())
}

func (s *Server) binaryRedBlueTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.BinaryRedBluePartition())
}

func (s *Server) lisTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.LISTrace())
}

func (s *Server) rowGravityTrace(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace.RowGravityTrace())
}

func (s *Server) flowTrace(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/traces/")
	result, ok := trace.FlowTrace(name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(result)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	nav := s.navigation(r.URL.Query().Get("tree"))
	data := homeData{
		navigationData: nav,
		Title:          "Fish Interview · 知识地图",
		Query:          query,
		Results:        s.search(query),
		Signals:        s.signalCounts(),
		TotalCards:     len(s.catalog.Cards),
	}
	s.render(w, "home", data)
}

func (s *Server) card(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/cards/")
	data, err := s.cardView(id, r.URL.Query().Get("tree"))
	if err != nil {
		http.Error(w, "render card", http.StatusInternalServerError)
		return
	}
	if data.Card.ID == "" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "card", data)
}

func (s *Server) cardPartial(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/partials/cards/")
	data, err := s.cardView(id, r.URL.Query().Get("tree"))
	if err != nil {
		http.Error(w, "render card", http.StatusInternalServerError)
		return
	}
	if data.Card.ID == "" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "card-content", data)
}

func (s *Server) cardView(id, requestedTree string) (cardData, error) {
	card, ok := s.catalog.Find(id)
	if !ok {
		return cardData{}, nil
	}
	body, err := s.renderMarkdown(card.Body)
	if err != nil {
		return cardData{}, err
	}
	related := make([]content.Card, 0, len(card.Links))
	for _, link := range card.Links {
		if target, ok := s.catalog.Find(link); ok {
			related = append(related, target)
		}
	}
	if requestedTree == "" {
		requestedTree = s.rootForCard(card.ID)
	}
	nav := s.navigation(requestedTree)
	nav.ActiveCardID = card.ID
	return cardData{
		navigationData: nav,
		Title:          card.Title + " · Fish Interview",
		Card:           card,
		Body:           template.HTML(body), // Generated from repository-owned Markdown.
		Related:        related,
		Signals:        card.ExamSignals,
		TraceURL:       traceURL(card.Trace),
	}, nil
}

func (s *Server) navigation(requestedID string) navigationData {
	options := make([]treeOption, 0, len(s.catalog.Roots))
	for _, root := range s.catalog.Roots {
		options = append(options, treeOption{ID: root.ID, Title: s.nodeTitle(root)})
	}
	selected := requestedID
	if !s.hasRoot(selected) && len(s.catalog.Roots) > 0 {
		selected = s.catalog.Roots[0].ID
	}
	roots := make([]treeViewNode, 0, len(s.catalog.Roots))
	for _, root := range s.catalog.Roots {
		roots = append(roots, s.viewNode(root))
	}
	payload, _ := json.Marshal(struct {
		Roots []treeViewNode `json:"roots"`
	}{Roots: roots})
	return navigationData{TreeID: selected, TreeOptions: options, TreeJSON: template.JS(payload)}
}

func (s *Server) viewNode(node content.TreeNode) treeViewNode {
	children := make([]treeViewNode, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, s.viewNode(child))
	}
	return treeViewNode{ID: node.ID, Title: s.nodeTitle(node), Card: node.Card, Children: children}
}

func (s *Server) nodeTitle(node content.TreeNode) string {
	if node.Title != "" {
		return node.Title
	}
	if card, ok := s.catalog.Find(node.Card); ok {
		return card.Title
	}
	return node.ID
}

func (s *Server) hasRoot(id string) bool {
	for _, root := range s.catalog.Roots {
		if root.ID == id {
			return true
		}
	}
	return false
}

func (s *Server) rootForCard(cardID string) string {
	for _, root := range s.catalog.Roots {
		if treeContains(root, cardID) {
			return root.ID
		}
	}
	return ""
}

func treeContains(node content.TreeNode, id string) bool {
	if node.ID == id || node.Card == id {
		return true
	}
	for _, child := range node.Children {
		if treeContains(child, id) {
			return true
		}
	}
	return false
}

func traceURL(name string) string {
	switch name {
	case "interval-scheduling":
		return "/api/traces/interval-scheduling"
	case "linear-dp":
		return "/api/traces/linear-dp"
	case "space-rolling":
		return "/api/traces/space-rolling"
	case "lcs-dp":
		return "/api/traces/lcs-dp"
	case "interval-dp":
		return "/api/traces/interval-dp"
	case "stock-dp":
		return "/api/traces/stock-dp"
	case "bitmask-dp":
		return "/api/traces/bitmask-dp"
	case "linked-list-rewire":
		return "/api/traces/linked-list-rewire"
	case "path-dp":
		return "/api/traces/path-dp"
	case "reverse-path-dp":
		return "/api/traces/reverse-path-dp"
	case "binary-red-blue":
		return "/api/traces/binary-red-blue"
	case "lis":
		return "/api/traces/lis"
	case "row-gravity":
		return "/api/traces/row-gravity"
	}
	if strings.HasPrefix(name, "flow-") {
		return "/api/traces/" + name
	}
	return ""
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "render page", http.StatusInternalServerError)
	}
}

func (s *Server) static(w http.ResponseWriter, _ *http.Request, name, contentType string) {
	data, err := assets.ReadFile(name)
	if err != nil {
		http.NotFound(w, nil)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(data)
}

func (s *Server) search(query string) []content.Card {
	all := s.catalog.SortedCards()
	if query == "" {
		return all
	}
	needle := strings.ToLower(query)
	matches := make([]content.Card, 0)
	for _, card := range all {
		haystack := strings.ToLower(strings.Join([]string{card.Title, card.Summary, strings.Join(card.Tags, " "), card.Body}, "\n"))
		if strings.Contains(haystack, needle) {
			matches = append(matches, card)
		}
	}
	return matches
}

func (s *Server) signalCounts() []signalCount {
	counts := map[string]int{}
	for _, card := range s.catalog.Cards {
		for _, signal := range card.ExamSignals {
			counts[signal.Company]++
		}
	}
	result := make([]signalCount, 0, len(counts))
	for company, count := range counts {
		result = append(result, signalCount{Company: company, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Company < result[j].Company
		}
		return result[i].Count > result[j].Count
	})
	return result
}

var wikilink = regexp.MustCompile(`\[\[([a-z0-9.-]+)\]\]`)
var codeBlock = regexp.MustCompile(`<pre([^>]*)>`)
var emphasisTrailingPunctuation = regexp.MustCompile(`\*\*([^*\n]+)([：:])\*\*`)
var latexBlock = regexp.MustCompile(`(?s)\$\$(.+?)\$\$`)
var latexSubscript = regexp.MustCompile(`_\{([^{}]+)\}`)
var latexSuperscript = regexp.MustCompile(`\^\{([^{}]+)\}`)
var latexText = regexp.MustCompile(`\\text\{([^{}]+)\}`)

func (s *Server) renderMarkdown(body string) (string, error) {
	body = emphasisTrailingPunctuation.ReplaceAllString(body, `**$1**$2`)
	formulas := make(map[string]string)
	body = latexBlock.ReplaceAllStringFunc(body, func(match string) string {
		key := fmt.Sprintf("FISHINTERVIEWMATH%d", len(formulas))
		formulas[key] = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "$$"), "$$"))
		return key
	})
	body = wikilink.ReplaceAllStringFunc(body, func(match string) string {
		id := wikilink.FindStringSubmatch(match)[1]
		if card, ok := s.catalog.Find(id); ok {
			return "[" + card.Title + "](/cards/" + id + ")"
		}
		return match
	})
	var rendered bytes.Buffer
	if err := s.markdown.Convert([]byte(body), &rendered); err != nil {
		return "", err
	}
	html := codeBlock.ReplaceAllString(rendered.String(), `<pre class="code-block"$1>`)
	for key, formula := range formulas {
		html = strings.ReplaceAll(html, "<p>"+key+"</p>", renderLatex(formula))
	}
	return html, nil
}

func renderLatex(formula string) string {
	formula = stdhtml.EscapeString(formula)
	formula = strings.ReplaceAll(formula, `\left`, "")
	formula = strings.ReplaceAll(formula, `\right`, "")
	formula = strings.ReplaceAll(formula, `\min`, "min")
	formula = strings.ReplaceAll(formula, `\max`, "max")
	formula = strings.ReplaceAll(formula, `\sum`, "∑")
	formula = strings.ReplaceAll(formula, `\infty`, "∞")
	formula = strings.ReplaceAll(formula, `\cup`, "∪")
	formula = strings.ReplaceAll(formula, `\le`, "≤")
	formula = strings.ReplaceAll(formula, `\ge`, "≥")
	formula = strings.ReplaceAll(formula, `\neq`, "≠")
	formula = strings.ReplaceAll(formula, `\;`, " ")
	formula = latexText.ReplaceAllString(formula, `<span class="math-text">$1</span>`)
	formula = latexSubscript.ReplaceAllString(formula, `<sub>$1</sub>`)
	formula = latexSuperscript.ReplaceAllString(formula, `<sup>$1</sup>`)
	formula = strings.ReplaceAll(formula, `\begin{cases}`, `<span class="math-cases">`)
	formula = strings.ReplaceAll(formula, `\end{cases}`, `</span>`)
	formula = strings.ReplaceAll(formula, `\\`, `<br>`)
	return `<div class="math-formula" aria-label="状态转移公式">` + formula + `</div>`
}

func tagClass(kind string) string {
	switch kind {
	case "algorithm-pattern":
		return "tag tag--algorithm"
	case "engineering":
		return "tag tag--engineering"
	default:
		return "tag"
	}
}
