package stool

import (
	"bufio"
	"bytes"
	"log"
	"testing"
)

func TestIndexEntireDirectory(t *testing.T) {
	buf := new(bytes.Buffer)
	fakeWriter := bufio.NewWriter(buf)

	indexer := &ViewIndexer{
		RootDir: "../test_fixtures/views2",
		Explainer: &VariableCollector{},
		ViewFinder: getViewFinder(),
		Writer: fakeWriter,
	}

	nodes := indexer.IndexViews("../test_fixtures/views2")

	root := nodes["root"]
	if root.Children[0] != "sibling" {
		log.Fatalf("expected to find child node with name of sibling, got %q" , root.Children[0])
	}

	if root.Children[1] != "dir1.view1" {
		log.Fatalf("expected to find child node with name of dir1.view1, got %q" , root.Children[1])
	}

	if len(root.Parents) != 0 {
		log.Fatalf("expected to not find parents, found %d parents" , len(root.Parents))
	}

	sibling := nodes["sibling"]
	if sibling.Children[0] != "dir1.view1" {
		log.Fatalf("expected to find child node with name of dir1.view1, got %q" , sibling.Children[0])
	}

	if sibling.Parents[0] != "root" {
		log.Fatalf("expected to find parent node with name of root, got %q" , sibling.Parents[0])
	}

	dir1view1 := nodes["dir1.view1"]
	if dir1view1.Parents[0] != "root" {
		log.Fatalf("expected to find parent node with name of dir1.view1, got %q" , sibling.Parents[0])
	}

	if dir1view1.Parents[1] != "sibling" {
		log.Fatalf("expected to find parent node with name of sibling, got %q" , sibling.Parents[1])
	}

	if len(dir1view1.Parents) != 2 {
		log.Fatalf("expected to find 2 parent nodes, got %q" , len(dir1view1.Parents))
	}

	dir1dir2view1 := nodes["dir1.dir2.view1"]
	if len(dir1dir2view1.Children) != 0 {
		log.Fatalf("expected to find 0 children nodes, got %d" , len(dir1dir2view1.Children))
	}

	if dir1dir2view1.Parents[0] != "dir1.view1" {
		log.Fatalf("expected to find dir1.view1 parent node, got %q" , dir1dir2view1.Parents[0])
	}

	if dir1dir2view1.Parents[1] != "dir1.view2" {
		log.Fatalf("expected to find dir1.view2 parent node, got %q" , dir1dir2view1.Parents[1])
	}
}

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

func newViewIndexer() *ViewIndexer {
	fakeWriter := bufio.NewWriter(new(bytes.Buffer))

	return &ViewIndexer{
		RootDir: "../test_fixtures/views2",
		Explainer: &VariableCollector{},
		ViewFinder: getViewFinder(),
		Writer: fakeWriter,
	}
}
