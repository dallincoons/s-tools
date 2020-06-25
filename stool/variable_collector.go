package stool

type ViewExplainer struct {
	ViewIndexer ViewIndexer
}

var parents = make(map[string]ViewNode)
var children = make(map[string]ViewNode)

type VariableCollection struct {
	Variables map[string]int
}

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

func (this *ViewExplainer) CollectChildrenFrom(viewName string) map[string]bool {
	children = map[string]ViewNode{}
	childrenNameSet := make(map[string]bool)

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectChildren(viewName, nodes)

	for _, p := range children {
		childrenNameSet[p.Name] = true
	}

	delete(childrenNameSet, viewName)

	return childrenNameSet
}

func (this *ViewExplainer) collectParents(viewName string, nodes map[string]ViewNode) {
	node, _  := nodes[viewName]

	parents[viewName] = node

	for _, parent := range node.Parents {
		this.collectParents(parent, nodes)
	}
}

func (this *ViewExplainer) collectChildren(viewName string, nodes map[string]ViewNode) {
	node, _  := nodes[viewName]

	children[viewName] = node

	for _, child := range node.Children {
		this.collectChildren(child, nodes)
	}
}

func (this *ViewExplainer) CollectVariablesFromParents(viewName string) map[string]int {
	collection := &VariableCollection{
		Variables: make(map[string]int),
	}

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectTreeVariables(viewName, nodes, collection)

	return collection.Variables
}

func (this *ViewExplainer) CollectVariablesFromChildren(viewName string) map[string]int {
	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	collection := &VariableCollection{
		Variables: make(map[string]int),
	}

	this.collectTreeVariablesDesc(viewName, nodes, collection)

	return collection.Variables
}

func (this *ViewExplainer) collectTreeVariables(viewName string, nodes map[string]ViewNode, collection *VariableCollection) {
	node, _  := nodes[viewName]

	for variable := range node.Variables {
		collection.Variables[variable]++
	}

	if len(node.Parents) > 0 {
		for _, parent := range node.Parents {
			this.collectTreeVariables(parent, nodes, collection)
		}
	}
}

func (this *ViewExplainer) collectTreeVariablesDesc(viewName string, nodes map[string]ViewNode, collection *VariableCollection) {
	node, _  := nodes[viewName]

	for variable := range node.Variables {
		collection.Variables[variable]++
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			this.collectTreeVariablesDesc(child, nodes, collection)
		}
	}
}
