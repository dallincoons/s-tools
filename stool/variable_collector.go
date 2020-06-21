package stool

type ViewExplainer struct {
	ViewIndexer ViewIndexer
}

var variables = make(map[string]int)
var parents = make(map[string]ViewNode)

func (this *ViewExplainer) CollectParentsFrom(viewName string) map[string]bool {
	parents = make(map[string]ViewNode)
	parentNameSet := make(map[string]bool)

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectParents(viewName, nodes)

	for _, p := range parents {
		parentNameSet[p.Name] = true
	}

	delete(parentNameSet, viewName)

	return parentNameSet
}

func (this *ViewExplainer) collectParents(viewName string, nodes map[string]ViewNode) {
	node, _  := nodes[viewName]

	parents[viewName] = node

	for _, parent := range node.Parents {
		this.collectParents(parent, nodes)
	}
}

func (this *ViewExplainer) CollectVariablesFromParents(viewName string) map[string]int {
	variables = make(map[string]int)

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectTreeVariables(viewName, nodes)

	return variables
}

func (this *ViewExplainer) CollectVariablesFromChildren(viewName string) map[string]int {
	variables = make(map[string]int)

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectTreeVariablesDesc(viewName, nodes)

	return variables
}

func (this *ViewExplainer) collectTreeVariables(viewName string, nodes map[string]ViewNode) {
	node, _  := nodes[viewName]

	for variable := range node.Variables {
		variables[variable]++
	}

	if len(node.Parents) > 0 {
		for _, parent := range node.Parents {
			this.collectTreeVariables(parent, nodes)
		}
	}
}

func (this *ViewExplainer) collectTreeVariablesDesc(viewName string, nodes map[string]ViewNode) {
	node, _  := nodes[viewName]

	for variable := range node.Variables {
		variables[variable]++
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			this.collectTreeVariablesDesc(child, nodes)
		}
	}
}
