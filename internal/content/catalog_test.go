package content

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDemoCatalog(t *testing.T) {
	catalog, err := Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got, want := len(catalog.Cards), 91; got != want {
		t.Fatalf("card count = %d, want %d", got, want)
	}
	if _, ok := catalog.Find("eng.kafka.rebalance"); !ok {
		t.Fatal("expected Kafka rebalance card")
	}
}

func TestAlgorithmPatternsStartWithLinkedExample(t *testing.T) {
	catalog, err := Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatal(err)
	}
	for id, card := range catalog.Cards {
		if card.Kind != "algorithm-pattern" {
			continue
		}
		if !strings.HasPrefix(card.Body, "## 例题\n") || !strings.Contains(card.Body, "https://leetcode.cn/problems/") {
			t.Fatalf("algorithm pattern %q must start with a linked example", id)
		}
	}
}

func TestAlgorithmPatternTreeCardsHaveTraces(t *testing.T) {
	catalog, err := Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatal(err)
	}
	var check func(TreeNode, bool)
	check = func(node TreeNode, underPatterns bool) {
		underPatterns = underPatterns || node.ID == "algo.patterns"
		if underPatterns && node.Card != "" && catalog.Cards[node.Card].Trace == "" {
			t.Fatalf("pattern tree card %q misses trace", node.Card)
		}
		for _, child := range node.Children {
			check(child, underPatterns)
		}
	}
	for _, root := range catalog.Roots {
		check(root, false)
	}
}

func TestParseCardRejectsMissingFrontmatter(t *testing.T) {
	if _, _, err := parseCard([]byte("# no frontmatter")); err == nil {
		t.Fatal("parseCard() error = nil, want error")
	}
}

func TestValidateRejectsPatternWithoutTrace(t *testing.T) {
	catalog := Catalog{
		Cards: map[string]Card{
			"algo.example": {
				Meta: Meta{ID: "algo.example", Kind: "algorithm-pattern", Title: "示例", Summary: "示例题型"},
			},
		},
	}
	if err := catalog.Validate(); err == nil || !strings.Contains(err.Error(), "misses trace") {
		t.Fatalf("Validate() error = %v, want missing trace", err)
	}
}
