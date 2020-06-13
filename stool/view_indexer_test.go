package stool

import (
	"log"
	"testing"
)

func TestFindAllIncludes(t *testing.T) {
	var siblingFound bool
	var view1Found bool

	view_indexer := newViewIndexer("../test_fixtures/views/root.blade.php")

	includes, _ := view_indexer.FindAllIncludes()

	for _, include := range includes {
		if (include == "sibling") {
			siblingFound = true
		}

		if (include == "dir1.view1") {
			view1Found = true
		}
	}

	if !siblingFound {
		log.Fatal("sibling include not found")
	}

	if !view1Found {
		log.Fatal("view1 include not found")
	}
}

func newViewIndexer(path string) *ViewIndexer {
	return &ViewIndexer{
		Path: path,
	}
}
