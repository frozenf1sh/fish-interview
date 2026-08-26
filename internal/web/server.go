package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/frozenf1sh/fish-interview/internal/content"
	"github.com/yuin/goldmark"
)

//go:embed templates/*.html static/*
var assets embed.FS

type Server struct {
	catalog   content.Catalog
	templates *template.Template
	markdown  goldmark.Markdown
}

type homeData struct {
	Title      string
	Query      string
	Roots      []content.TreeNode
	Results    []content.Card
	Signals    []signalCount
	TotalCards int
}

type cardData struct {
	Title   string
	Card    content.Card
	Body    template.HTML
	Related []content.Card
	Signals []content.ExamSignal
	Roots   []content.TreeNode
}

type signalCount struct {
	Company string
	Count   int
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
	return &Server{catalog: catalog, templates: templates, markdown: goldmark.New()}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		s.home(w, r)
	case strings.HasPrefix(r.URL.Path, "/cards/"):
		s.card(w, r)
	case r.URL.Path == "/static/app.css":
		s.static(w, r, "static/app.css", "text/css; charset=utf-8")
	case r.URL.Path == "/static/app.js":
		s.static(w, r, "static/app.js", "text/javascript; charset=utf-8")
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	data := homeData{
		Title:      "Fish Interview · 知识地图",
		Query:      query,
		Roots:      s.catalog.Roots,
		Results:    s.search(query),
		Signals:    s.signalCounts(),
		TotalCards: len(s.catalog.Cards),
	}
	s.render(w, "home", data)
}

func (s *Server) card(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/cards/")
	card, ok := s.catalog.Find(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	body, err := s.renderMarkdown(card.Body)
	if err != nil {
		http.Error(w, "render card", http.StatusInternalServerError)
		return
	}
	related := make([]content.Card, 0, len(card.Links))
	for _, link := range card.Links {
		if target, ok := s.catalog.Find(link); ok {
			related = append(related, target)
		}
	}
	s.render(w, "card", cardData{
		Title:   card.Title + " · Fish Interview",
		Card:    card,
		Body:    template.HTML(body), // Generated from repository-owned Markdown.
		Related: related,
		Signals: card.ExamSignals,
		Roots:   s.catalog.Roots,
	})
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

func (s *Server) renderMarkdown(body string) (string, error) {
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
	return rendered.String(), nil
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
