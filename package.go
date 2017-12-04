package antho

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/blang/semver"
	"github.com/google/go-jsonnet"
	"github.com/pkg/errors"
)

// JSonnetPath is the env variable to set the import paths.
const JSonnetPath = "JSONNET_PATH"

// PackageFile is the default package configuration filename.
// Note: This is your `package.json` or `Gemfile` equivalent.
const PackageFile = "antho.json"

// IndexFile is the default index file to read for a library.
const IndexFile = "main.libsonnet"

// ExcludeMatches are the list of file matches to ignore when building a package.
var ExcludeMatches = []string{
	"vendor",
	"vendor/*",
}

// FromFile takes a given directory path and attempts to parse the package.
func FromFile(dpath string) (*Package, error) {
	pfile := path.Join(dpath, PackageFile)
	pdata, err := ioutil.ReadFile(pfile)
	if err != nil {
		return nil, err
	}

	return Parse(dpath, pdata)
}

// Parse will take in a given package json and return the go struct.
func Parse(dpath string, contents []byte) (*Package, error) {
	var pkg Package

	err := json.Unmarshal(contents, &pkg)
	if err != nil {
		return nil, err
	}

	pkg.Path = dpath
	return &pkg, err
}

// Package represents a jsonnet library.
type Package struct {
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	Dependencies []Package `json:"dependencies"`

	Path string
}

// Semver returns the semantic version for the package.
func (p *Package) Semver() (semver.Version, error) {
	return semver.Make(p.Version)
}

// Files returns the list of files inside of the package.
// Note: This excludes the directory path itself.
func (p *Package) Files() ([]string, error) {
	fileList := []string{}

	err := filepath.Walk(p.Path, func(fpath string, f os.FileInfo, e error) error {
		excluded, e := p.isExcludedFile(fpath)
		if e != nil {
			return e
		}

		if fpath != p.Path && !excluded {
			rel, e := p.relPath(fpath)
			if e != nil {
				return e
			}

			fileList = append(fileList, rel)
		}
		return e
	})
	if err != nil {
		return []string{}, err
	}

	sort.Strings(fileList)
	return fileList, err
}

func (p *Package) relPath(filename string) (string, error) {
	return filepath.Rel(p.Path, filename)
}

func (p *Package) absPath(relpath string) (string, error) {
	return filepath.Abs(path.Join(p.Path, relpath))
}

func (p *Package) isExcludedFile(filename string) (bool, error) {
	for _, m := range ExcludeMatches {
		rel, err := p.relPath(filename)
		if err != nil {
			return true, errors.Wrapf(err, "unable to retrieve relPath for `%s`", filename)
		}

		matched, err := filepath.Match(m, rel)
		if err != nil {
			return true, err
		}

		if matched {
			return matched, nil
		}
	}

	return false, nil
}

// Build will generate a tarball of the given package directory.
func (p *Package) Build(w io.Writer) error {
	tw := tar.NewWriter(w)
	files, err := p.Files()
	if err != nil {
		tw.Close()
		return err
	}

	for _, f := range files {
		fi, err := os.Stat(path.Join(p.Path, f))
		if err != nil {
			tw.Close()
			return err
		}

		hdr := &tar.Header{
			Name: fi.Name(),
			Mode: int64(fi.Mode()),
			Size: fi.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			tw.Close()
			return err
		}

		absPath, err := p.absPath(f)
		if err != nil {
			tw.Close()
			return err
		}

		input, err := os.Open(absPath)
		if err != nil {
			tw.Close()
			return err
		}
		defer input.Close()

		if _, err := io.Copy(tw, input); err != nil {
			tw.Close()
			return err
		}
	}

	return tw.Close()
}

// TarballName is the name of the tarball package.
func (p *Package) TarballName() string {
	return fmt.Sprintf("%s-%s.tar.gz", p.Name, p.Version)
}

// WriteTarball is a helper to build and write the output to a file.
func (p *Package) WriteTarball(targetDir string) error {
	targetFile := path.Join(targetDir, p.TarballName())
	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	err = p.Build(file)
	if err != nil {
		file.Close()
		return err
	}

	return file.Close()
}

// Graph returns the dependency tree of packages.
func (p *Package) Graph() *Graph {
	return GenerateGraph(p)
}

// JPath tries to return the jsonnet include path for a given package.
// This is used for validating a package, running tests,
// and potentially will/can be used for integrating with editors.
func (p *Package) JPath() []string {
	env := []string{p.Path}
	if jPath := os.Getenv(JSonnetPath); jPath != "" {
		env = append(env, jPath)
	}

	return env
}

// Validate will check test whether or not a package is ready to be built.
func (p *Package) Validate() error {
	index := path.Join(p.Path, IndexFile)
	pdata, err := ioutil.ReadFile(index)
	if err != nil {
		err = fmt.Errorf("validation error: %s", err)
		return err
	}

	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: p.JPath(),
	})

	_, err = vm.EvaluateSnippet(IndexFile, string(pdata))
	if err != nil {
		err = fmt.Errorf("validation error: %s", err)
	}
	return err
}
