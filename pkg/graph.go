package antho

// Graph represents a dependency graph.
type Graph struct {
	Package *Package
}

func (g *Graph) String() string {
	return "awesome"
}

// GenerateGraph takes a given package and returns
// the dependency graph for it.
func GenerateGraph(pkg *Package) *Graph {
	g := &Graph{
		Package: pkg,
	}

	return g
}
