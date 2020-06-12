package cmd

import (
	"log"
	"testing"
	"surgio-tools/stool"
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

func getViewExplainer(rootPath string) *stool.ViewExplainer {
	return &stool.ViewExplainer{
		RootPath: rootPath,
	}
}
