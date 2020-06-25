package stool

type ViewExplainer struct {
	ViewIndexer ViewIndexer
}

type VariableCollection struct {
	Variables map[string]int
}

type ParentCollection struct {
	Parents map[string]ViewNode
}

type ChildrenCollection struct {
	Children map[string]ViewNode
}

func (this *ViewExplainer) CollectParentsFrom(viewName string) map[string]bool {
	parentCollection := &ParentCollection{
		Parents:make(map[string]ViewNode),
	}
	parentNameSet := make(map[string]bool)

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectParents(viewName, nodes, parentCollection)

	for _, p := range parentCollection.Parents {
		parentNameSet[p.Name] = true
	}

	delete(parentNameSet, viewName)

	return parentNameSet
}

func (this *ViewExplainer) CollectChildrenFrom(viewName string) map[string]bool {
	childrenCollection := &ChildrenCollection{Children: make(map[string]ViewNode)}
	childrenNameSet := make(map[string]bool)

	nodes := this.ViewIndexer.IndexViews(this.ViewIndexer.RootDir)

	this.collectChildren(viewName, nodes, childrenCollection)

	for _, p := range childrenCollection.Children {
		childrenNameSet[p.Name] = true
	}

	delete(childrenNameSet, viewName)

	return childrenNameSet
}

func (this *ViewExplainer) collectParents(viewName string, nodes map[string]ViewNode, collection *ParentCollection) {
	node, _  := nodes[viewName]

	collection.Parents[viewName] = node

	for _, parent := range node.Parents {
		this.collectParents(parent, nodes, collection)
	}
}

func (this *ViewExplainer) collectChildren(viewName string, nodes map[string]ViewNode, collection *ChildrenCollection) {
	node, _  := nodes[viewName]

	collection.Children[viewName] = node

	for _, child := range node.Children {
		this.collectChildren(child, nodes, collection)
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
