package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frozenf1sh/fish-interview/internal/trace"
)

func main() {
	failures := trace.ValidateAllPlayerContracts()
	for name, err := range failures {
		fmt.Fprintf(os.Stderr, "trace %s: %v\n", name, err)
	}
	for _, directory := range []string{"content", "internal"} {
		if err := scanVisibleSource(directory); err != nil {
			fmt.Fprintln(os.Stderr, err)
			failures[directory] = err
		}
	}
	if len(failures) > 0 {
		os.Exit(1)
	}
	fmt.Println("tracecheck: all traces renderable and user-visible source has no HTML operator entities")
}

func scanVisibleSource(directory string) error {
	return filepath.WalkDir(directory, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lessEntity, greaterEntity := "&"+"lt;", "&"+"gt;"
		if strings.Contains(string(data), lessEntity) || strings.Contains(string(data), greaterEntity) {
			return fmt.Errorf("%s contains HTML-escaped comparison operator", path)
		}
		return nil
	})
}
