package stool

import (
	"io/ioutil"
	"regexp"
)

type ViewIndexer struct {
	Explainer *ViewExplainer
	ViewFinder *ViewFinder
}

type ViewNode struct {
	Name string
	children []*ViewNode
	parent *ViewNode
	variables map[string]int
}

var tree *ViewNode

func (this *ViewIndexer) index(view_name string) *ViewNode {
	path := this.ViewFinder.GetFilePath(view_name)
	includes, err := this.FindAllIncludes(path)

	if err == nil {
		variables, _ := this.Explainer.GetAllVariablesFrom(path)

		tree = &ViewNode{
			Name: view_name,
			variables: variables,
		}

		this.addChildren(tree, includes)
	}

	return tree
}

func (this *ViewIndexer) addChildren(parentNode *ViewNode, includes []string) {
	var path string
	var child *ViewNode

	for _, view_name := range includes {
		path = this.ViewFinder.GetFilePath(view_name)

		variables, _ := this.Explainer.GetAllVariablesFrom(path)

		child = &ViewNode {
			Name: view_name,
			parent: parentNode,
			variables: variables,
		}

		parentNode.children = append(parentNode.children, child)

		includes, _ := this.FindAllIncludes(this.ViewFinder.GetFilePath(view_name))

		this.addChildren(child, includes)
	}
}

func (this *ViewIndexer) FindAllIncludes(path string) ([]string, error) {
	var viewNames []string

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return this.getViewNames(contents, viewNames), nil
}

func (this *ViewIndexer) getViewNames(contents []byte, viewNames []string) []string {
	re := regexp.MustCompile("@include\\('(.+)'\\)")

	results := re.FindAllStringSubmatch(string(contents), -1)

	for _, result := range results {
		viewNames = append(viewNames, result[1])
	}

	return viewNames
}
