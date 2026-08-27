package content

import (
	"path/filepath"
	"testing"
)

func TestLoadDemoCatalog(t *testing.T) {
	catalog, err := Load(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got, want := len(catalog.Cards), 11; got != want {
		t.Fatalf("card count = %d, want %d", got, want)
	}
	if _, ok := catalog.Find("eng.kafka.rebalance"); !ok {
		t.Fatal("expected Kafka rebalance card")
	}
}

func TestParseCardRejectsMissingFrontmatter(t *testing.T) {
	if _, _, err := parseCard([]byte("# no frontmatter")); err == nil {
		t.Fatal("parseCard() error = nil, want error")
	}
}
