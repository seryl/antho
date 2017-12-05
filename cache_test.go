package antho

import (
	"path"
	"testing"
)

func TestPackageMap(t *testing.T) {
	cachePath := path.Join(fixturesDirectory, "cache")
	cache := &Cache{Path: cachePath}

	pkgs, err := cache.Packages()
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) != 2 {
		t.Error("Package quantity is incorrect")
	}
}
