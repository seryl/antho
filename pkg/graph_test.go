package antho

import (
	"testing"
)

func TestGenerateDependencyGraph(t *testing.T) {
	pkg, err := fixturePackage("cache/github.com/seryl/examplewithdep-0.0.1")
	if err != nil {
		t.Error(err)
	}

	graph := pkg.Graph()
	t.Log(graph)
}
