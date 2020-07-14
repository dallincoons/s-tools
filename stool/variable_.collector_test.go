package stool

import (
	"log"
	"testing"
)

func TestGetAllParentsInTree(t *testing.T) {
	explainer := newExplainer()

	parents := explainer.CollectParentsFrom("dir1.dir2.view1")

	if !parents["root"] {
		log.Fatal("root not found as a parent")
	}

	if !parents["sibling"] {
		log.Fatal("sibling not found as a parent")
	}

	if !parents["dir1.view1"] {
		log.Fatal("dir1.view1 not found as a parent")
	}

	if !parents["dir1.view2"] {
		log.Fatal("dir1.view2 not found as a parent")
	}

	if parents["dir1.view3"] {
		log.Fatal("dir1.view3 unexpectantly found as a parent")
	}
}

func TestGetAllChildrenInTree(t *testing.T) {
	explainer := newExplainer()

	children := explainer.CollectChildrenFrom("root")

	if !children["dir1.dir2.view1"] {
		log.Fatal("dir1.view2.view1 not found as a child")
	}

	if !children["dir1.view1"] {
		log.Fatal("dir1.view1 not found as a child")
	}

	if !children["sibling"] {
		log.Fatal("sibling not found as a child")
	}
}

func TestGetAllVariablesFromParentsInTree(t *testing.T) {
	explainer := newExplainer()

	variables := explainer.CollectVariablesFromParents("dir1.dir2.view1")

	_, jimmy_found := variables["$jimmy"]

	if !jimmy_found {
		log.Fatal("jimmy variable not found")
	}
}

func TestGetAllVariablesFromChildrenInTree(t *testing.T) {
	explainer := newExplainer()

	variables := explainer.CollectVariablesFromChildren("cousin")

	torn_count, torn_found := variables["$torn"]

	if !torn_found {
		log.Fatal("torn variable not found")
	}

	if torn_count != 1 {
		log.Fatalf("torn variable should be found 1 time, found %d times", torn_count)
	}

	blondie_count, blondie_found := variables["$blondie"]

	if !blondie_found {
		log.Fatal("blondie variable not found")
	}

	if blondie_count != 1 {
		log.Fatalf("blondie variable should be found 1 time, found %d times", blondie_count)
	}

	kinks_count, kinks_found := variables["$kinks"]

	if !kinks_found {
		log.Fatal("kinks variable not found")
	}

	if kinks_count != 1 {
		log.Fatalf("kinks variable should be found 1 time, found %d times", kinks_count)
	}

	if len(variables) != 3 {
		log.Fatal("expected only 3 variables to exist in the cousin tree")
	}
}

func newExplainer() *ViewExplainer {
	return &ViewExplainer{
		ViewIndexer: ViewIndexer{
			RootDir:    "../test_fixtures/views2/resources/views",
			Explainer:  &VariableCollector{},
			ViewFinder: &ViewFinder{},
			Writer:     nil,
		},
	}
}

