package stool

import (
	"log"
	"testing"
)

func TestFindAllIncludes(t *testing.T) {
	var siblingFound bool
	var view1Found bool

	view_indexer := newViewIndexer()

	includes, _ := view_indexer.FindAllIncludes("../test_fixtures/views/root.blade.php")

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

func TestBuildViewTree(t *testing.T) {
	view_indexer := newViewIndexer()

	tree := view_indexer.Index("root")

	if (tree.Name != "root") {
		log.Fatal("expected to find root node with name of root")
	}

	if (tree.Children[0].Name != "account.sub-header") {
		log.Fatal("expected to find child node with name of sibling")
	}

	if (tree.Children[1].Name != "sibling") {
		log.Fatal("expected to find child node with name of sibling")
	}

	if (tree.Children[2].Name != "dir1.view1") {
		log.Fatal("expected to find child node with name of dir1.view1")
	}

	if (tree.Children[2].Children[1].Name != "ops._nav") {
		log.Fatal("expected to find child node with name of ops._nav")
	}

	if (tree.Children[2].Children[0].Name != "dir2.view2") {
		log.Fatal("expected to find child node with name of dir2.view2")
	}
}

func newViewIndexer() *ViewIndexer {
	return &ViewIndexer{
		&ViewExplainer{},
		getViewFinder(),
	}
}
