package stool

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func (this *ViewIndexer) IndexViews(root string) map[string]ViewNode {
	blade := this.indexBladeViews(root)

	var p string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(),".blade.php") {
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
	child_already_exists := false
	parent_already_exists := false

	includes, _ := this.FindAllIncludes(filepath.Join(this.RootDir, path))

	for _, include := range includes {
		childNode, found := nodes[include]

		if found {
			for _, c := range node.Children {
				child_already_exists = (c == include)
				if child_already_exists {
					break
				}
			}

			for _, c := range childNode.Parents {
				parent_already_exists = (c == node.Name)
				if parent_already_exists {
					break
				}
			}

			if !child_already_exists {
				node.Children = append(node.Children, include)
			}

			if !parent_already_exists {
				childNode.Parents = append(childNode.Parents, node.Name)
			}

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

func (this *ViewIndexer) indexBladeViews(root string) *Blade {
	blade := &Blade{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(),".blade.php") {
			return nil
		}

		rel, err := filepath.Rel(root, path)

		if err != nil {
			return err
		}

		blade.AddPath(rel)

		this.storeNode(blade, rel)

		return nil
	})

	return blade
}

func (this *ViewIndexer) storeNode(blade *Blade, rel string) {
	viewName, _ := blade.GetName(rel)

	node := this.getNode(viewName, *blade)

	nodes[viewName] = *node
}

func (this *ViewIndexer) IndexView(view_name string) []*ViewNode {
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
