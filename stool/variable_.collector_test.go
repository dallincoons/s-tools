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

func TestGetAllVariablesInTree(t *testing.T) {
	explainer := newExplainer()

	variables := explainer.CollectVariablesFrom("dir1.dir2.view1")

	_, jimmy_found := variables["$jimmy"]

	if !jimmy_found {
		log.Fatal("jimmy variable not found")
	}
}

func newExplainer() *ViewExplainer {
	return &ViewExplainer{
		ViewIndexer: ViewIndexer{
			RootDir:    "../test_fixtures/views2",
			Explainer:  &VariableCollector{},
			ViewFinder: &ViewFinder{},
			Writer:     nil,
		},
	}
}

