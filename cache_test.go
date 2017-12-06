package antho

import (
	"path"
	"testing"
)

func TestPackageMap(t *testing.T) {
	cachePath := path.Join(fixturesDirectory, "cache")
	cache := &Cache{Path: cachePath}

	expected := []struct {
		name    string
		origin  string
		version string
	}{
		{
			name:    "example",
			origin:  "github.com/seryl",
			version: "0.0.1",
		},
		{
			name:    "exampledep",
			origin:  "github.com/seryl",
			version: "0.0.1",
		},
	}

	pkgs, err := cache.Packages()
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) != len(expected) {
		t.Error("Package quantity is incorrect")
	}

	for _, pkg := range expected {
		pkgPath := path.Join(pkg.origin, pkg.name)
		if pkgs[pkgPath] == nil {
			t.Errorf("Package path does not exist: `%s`", pkgPath)
		}

		if pkgs[pkgPath][pkg.version] == nil {
			t.Errorf("Package version exist: `%s`, `%s`", pkgPath, pkg.version)
		}

		p := pkgs[pkgPath][pkg.version]
		if p.Name != pkg.name {
			t.Errorf("Names do not match. Expected: `%s`, Received: `%s`",
				pkg.name, p.Name)
		}

		if p.Version != pkg.version {
			t.Errorf("Versions do not match. Expected: `%s`, Received: `%s`",
				pkg.version, p.Version)
		}

		if p.Origin != pkg.origin {
			t.Errorf("Origins do not match. Expected: `%s`, Received: `%s`",
				pkg.origin, p.Origin)
		}
	}
}
