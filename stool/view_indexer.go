package stool

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var index = make(map[string]*ViewNode)
var nodes = make(map[string]ViewNode)

type ViewIndexer struct {
	RootDir string
	Explainer *VariableCollector
	ViewFinder *ViewFinder
	Writer *bufio.Writer
}

type ViewNode struct {
	Name string
	Children []string
	Parents []string
	Variables map[string]int
}

var root *ViewNode
var rx = regexp.MustCompile("(.+).blade.php")

func (this *ViewIndexer) IndexViews(root string) map[string]ViewNode {
	blade := this.storeBladeViewNames(root)

	var p string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		p, _ = filepath.Rel(root, path)
		viewName, _ := blade.GetName(p)

		node := this.getNode(viewName, *blade)

		nodes[viewName] = *node

		return nil
	})

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		p, _ = filepath.Rel(root, path)
		viewName, _ := blade.GetName(p)

		node := this.getNode(viewName, *blade)

		this.addNodeChildren(node, blade)

		return nil
	})

	return nodes
}

func (this *ViewIndexer) addNodeChildren(node *ViewNode, blade *Blade) {
	path, _ := blade.GetPath(node.Name)

	includes, _ := this.FindAllIncludes(this.RootDir + string(filepath.Separator) + path)

	for _, include := range includes {
		childNode, found := nodes[include]

		if found {
			childNode.Parents = append(childNode.Parents, node.Name)
			node.Children = append(node.Children, include)

			nodes[include] = childNode
			nodes[node.Name] = *node
		}
	}
}

func (this *ViewIndexer) getNode(name string, blade Blade) *ViewNode {
	node, found := nodes[name]
	path, _ := blade.GetPath(name)
	variables, _ := this.Explainer.GetAllVariablesFrom(this.RootDir + string(filepath.Separator) + path)

	if (!found) {
		return &ViewNode{
			Name: name,
			Children: []string{},
			Parents: []string{},
			Variables: variables,
		}
	}

	return &node
}

func (this *ViewIndexer) storeBladeViewNames(root string) *Blade {
	blade := &Blade{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)

		if err != nil {
			return err
		}

		blade.AddPath(rel)

		return nil
	})

	return blade
}

func (this *ViewIndexer) IndexView(view_name string) []*ViewNode {
	fmt.Println("test")
	path := this.ViewFinder.GetFilePath(view_name)
	includes, err := this.FindAllIncludes(path)
	var flatTree []*ViewNode

	if err == nil {
		variables, _ := this.Explainer.GetAllVariablesFrom(path)

		root = &ViewNode{
			Name: view_name,
			Variables: variables,
		}

		flatTree = this.addChildren(root, includes)
		flatTree = append(flatTree, root)
	}

	return flatTree
}

func (this *ViewIndexer) addChildren(parentNode *ViewNode, includes []string) []*ViewNode {
	var path string
	var child *ViewNode
	var flatTree []*ViewNode

	for _, view_name := range includes {
		path = this.ViewFinder.GetFilePath(view_name)
		variables, _ := this.Explainer.GetAllVariablesFrom(path)

		if _, found := index[view_name]; found == true {
			index[view_name].Parents = append(index[view_name].Parents, parentNode.Name)
			continue
		}

		child = &ViewNode {
			Name: view_name,
			Parents: []string{parentNode.Name},
			Variables: variables,
		}

		parentNode.Children = append(parentNode.Children, child.Name)
		flatTree = append(flatTree, child)
		index[child.Name] = child

		includes, _ := this.FindAllIncludes(this.ViewFinder.GetFilePath(view_name))

		this.addChildren(child, includes)
	}

	return flatTree
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
	re := regexp.MustCompile("@include\\('(.+)\\'")

	results := re.FindAllStringSubmatch(string(contents), -1)

	for _, result := range results {
		viewNames = append(viewNames, result[1])
	}

	return viewNames
}
