package stool

import (
	"log"
	"testing"
)

func TestGetErrorIfViewDoesNotExist(t *testing.T) {
	explainer := getViewExplainer()

	_, error := explainer.GetAllVariablesFrom("dudnt exist")

	if error == nil {
		log.Fatal("Expecting an error due to a reference to non-existent view")
	}
}

func TestGetAllVariablesInAView(t *testing.T) {
	explainer := getViewExplainer()

	variables, _ := explainer.GetAllVariablesFrom("../test_fixtures/views/root.blade.php")

	_, found := variables["$title"]

	if !found {
		log.Fatal("$title variable not found")
	}
}

func getViewExplainer() *VariableCollector {
	return &VariableCollector{}
}
