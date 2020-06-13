package stool

import (
	"log"
	"testing"
)

func TestGetErrorIfViewDoesNotExist(t *testing.T) {
	explainer := getViewExplainer("../test_fixtures/views");

	_, error := explainer.GetAllVariablesFrom("dudnt exist")

	if error == nil {
		log.Fatal("Expecting an error due to a reference to non-existent view")
	}
}

func TestGetAllVariablesInAView(t *testing.T) {
	explainer := getViewExplainer("../test_fixtures/views");

	variables, _ := explainer.GetAllVariablesFrom("root")

	_, found := variables["$title"]

	if !found {
		log.Fatal("$title variable not found")
	}
}

func getViewExplainer(rootPath string) *ViewExplainer {
	return &ViewExplainer{
		RootPath: rootPath,
	}
}
