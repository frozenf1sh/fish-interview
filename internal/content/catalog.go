package content

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ExamSignal is time-bounded external evidence, not a claim that a company owns a topic.
type ExamSignal struct {
	Company    string `yaml:"company"`
	Year       int    `yaml:"year"`
	Role       string `yaml:"role"`
	Confidence string `yaml:"confidence"`
	Source     string `yaml:"source"`
}

type Meta struct {
	ID          string       `yaml:"id"`
	Kind        string       `yaml:"kind"`
	Title       string       `yaml:"title"`
	Summary     string       `yaml:"summary"`
	Parents     []string     `yaml:"parents"`
	Tags        []string     `yaml:"tags"`
	Links       []string     `yaml:"links"`
	Trace       string       `yaml:"trace"`
	ExamSignals []ExamSignal `yaml:"exam_signals"`
}

type Card struct {
	Meta
	Body string
	Path string
}

type TreeNode struct {
	ID       string     `yaml:"id"`
	Title    string     `yaml:"title"`
	Card     string     `yaml:"card"`
	Children []TreeNode `yaml:"children"`
}

type treeFile struct {
	Roots []TreeNode `yaml:"roots"`
}

type Catalog struct {
	Cards map[string]Card
	Roots []TreeNode
}

func Load(contentDir string) (Catalog, error) {
	cards, err := loadCards(filepath.Join(contentDir, "cards"))
	if err != nil {
		return Catalog{}, err
	}

	data, err := os.ReadFile(filepath.Join(contentDir, "taxonomy", "tree.yaml"))
	if err != nil {
		return Catalog{}, fmt.Errorf("read taxonomy: %w", err)
	}
	var tree treeFile
	if err := yaml.Unmarshal(data, &tree); err != nil {
		return Catalog{}, fmt.Errorf("parse taxonomy: %w", err)
	}

	catalog := Catalog{Cards: cards, Roots: tree.Roots}
	if err := catalog.Validate(); err != nil {
		return Catalog{}, err
	}
	return catalog, nil
}

func loadCards(dir string) (map[string]Card, error) {
	cards := make(map[string]Card)
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		meta, body, err := parseCard(data)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if _, exists := cards[meta.ID]; exists {
			return fmt.Errorf("duplicate card ID %q", meta.ID)
		}
		cards[meta.ID] = Card{Meta: meta, Body: body, Path: path}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load cards: %w", err)
	}
	return cards, nil
}

func parseCard(data []byte) (Meta, string, error) {
	text := string(data)
	if !strings.HasPrefix(text, "---\n") {
		return Meta{}, "", errors.New("frontmatter must start with ---")
	}
	rest := strings.TrimPrefix(text, "---\n")
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return Meta{}, "", errors.New("frontmatter must end with ---")
	}
	var meta Meta
	if err := yaml.Unmarshal([]byte(rest[:end]), &meta); err != nil {
		return Meta{}, "", err
	}
	return meta, strings.TrimSpace(rest[end+5:]), nil
}

func (c Catalog) Validate() error {
	if len(c.Cards) == 0 {
		return errors.New("catalog has no cards")
	}
	nodeIDs := make(map[string]bool)
	for _, root := range c.Roots {
		if err := validateNode(root, c.Cards, nodeIDs); err != nil {
			return err
		}
	}
	for id, card := range c.Cards {
		if card.ID == "" || card.Kind == "" || card.Title == "" || card.Summary == "" {
			return fmt.Errorf("card %q misses required metadata", id)
		}
		for _, parent := range card.Parents {
			if !nodeIDs[parent] {
				return fmt.Errorf("card %q refers to unknown parent %q", id, parent)
			}
		}
		for _, link := range card.Links {
			if _, ok := c.Cards[link]; !ok {
				return fmt.Errorf("card %q links to unknown card %q", id, link)
			}
		}
		for _, signal := range card.ExamSignals {
			if signal.Company == "" || signal.Year == 0 || signal.Role == "" || signal.Confidence == "" || signal.Source == "" {
				return fmt.Errorf("card %q has incomplete exam signal", id)
			}
		}
	}
	return nil
}

func validateNode(node TreeNode, cards map[string]Card, seen map[string]bool) error {
	if node.ID == "" {
		return errors.New("taxonomy node misses ID")
	}
	if seen[node.ID] {
		return fmt.Errorf("duplicate taxonomy node ID %q", node.ID)
	}
	seen[node.ID] = true
	if node.Title == "" && node.Card == "" {
		return fmt.Errorf("taxonomy node %q misses title", node.ID)
	}
	if node.Card != "" {
		if node.ID != node.Card {
			return fmt.Errorf("node %q must match its card ID %q", node.ID, node.Card)
		}
		if _, ok := cards[node.Card]; !ok {
			return fmt.Errorf("taxonomy node %q refers to missing card", node.ID)
		}
	}
	for _, child := range node.Children {
		if err := validateNode(child, cards, seen); err != nil {
			return err
		}
	}
	return nil
}

func (c Catalog) SortedCards() []Card {
	cards := make([]Card, 0, len(c.Cards))
	for _, card := range c.Cards {
		cards = append(cards, card)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].Title < cards[j].Title })
	return cards
}

func (c Catalog) Find(id string) (Card, bool) {
	card, ok := c.Cards[id]
	return card, ok
}
