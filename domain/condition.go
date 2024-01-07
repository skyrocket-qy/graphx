package domain

type SearchCondition struct {
	In Compare `json:"in"`
}

type Compare struct {
	Namespaces []string `json:"namespaces"`
	Names      []string `json:"names"`
	Relations  []string `json:"relations"`
}

func (c *SearchCondition) ShouldStop(node Node) bool {
	if len(c.In.Namespaces) == 0 && len(c.In.Names) == 0 && len(c.In.Relations) == 0 {
		// means no specific condition, never stop
		return false
	}
	for _, namespace := range c.In.Namespaces {
		if node.Namespace == namespace {
			return false
		}
	}
	for _, name := range c.In.Names {
		if node.Name == name {
			return false
		}
	}
	for _, relation := range c.In.Relations {
		if node.Relation == relation {
			return false
		}
	}
	return true
}

type CollectCondition struct {
	In Compare `json:"in"`
}

func (c *CollectCondition) ShouldCollect(node Node) bool {
	if len(c.In.Namespaces) == 0 && len(c.In.Names) == 0 && len(c.In.Relations) == 0 {
		// means no specific conditions, collect all nodes
		return true
	}
	for _, namespace := range c.In.Namespaces {
		if node.Namespace == namespace {
			return true
		}
	}
	for _, name := range c.In.Names {
		if node.Name == name {
			return true
		}
	}
	for _, relation := range c.In.Relations {
		if node.Relation == relation {
			return true
		}
	}
	return false
}
