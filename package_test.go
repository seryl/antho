package antho

import (
	"archive/tar"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"
	"sort"
	"testing"
)

func TestPackageFailure(t *testing.T) {
	_, err := fixturePackage("cache/github.com/seryl/does-not-exist-0.0.1")
	if err == nil {
		t.Error("Package parsing is not failing with non-existant package")
	}
}

func TestPackageFiles(t *testing.T) {
	expectedContents := []string{
		"antho.json",
		"main.libsonnet",
	}
	sort.Strings(expectedContents)

	pkg, err := fixturePackage("cache/github.com/seryl/example-0.0.1")
	if err != nil {
		t.Error(err)
	}

	if pkg == nil {
		t.Error("Package is nil")
	}

	fList, err := pkg.Files()
	if err != nil {
		t.Error(err)
	}
	sort.Strings(fList)

	for i := 0; i < len(expectedContents); i++ {
		if expectedContents[i] != fList[i] {
			t.Errorf("Contents do not match. Expected: `%s`, Received: `%s`",
				expectedContents[i], fList[i])
		}
	}

	if len(fList) != len(expectedContents) {
		t.Errorf("Sizes do not match. Expected: `%s`, Received: `%s`",
			expectedContents, fList)
	}
}

func TestValidation(t *testing.T) {
	pkg, err := fixturePackage("cache/github.com/seryl/example-0.0.1")
	if err != nil {
		t.Error(err)
	}

	err = pkg.Validate()
	if err != nil {
		t.Error(err)
	}
}

func TestJPath(t *testing.T) {
	examplePath := "/my/example"
	pkg, err := fixturePackage("cache/github.com/seryl/example-0.0.1")
	if err != nil {
		t.Error(err)
	}

	jpath, err := pkg.JPath()
	if err != nil {
		t.Error(err)
	}

	index := sort.SearchStrings(jpath, pkg.Path)
	if index == len(jpath) {
		t.Errorf("Unable to find package path in JSonnet Path")
		t.Errorf("Current %s: %s", JSonnetPath, jpath)
	}

	err = os.Setenv(JSonnetPath, examplePath)
	if err != nil {
		t.Errorf("Unable to set Jsonnet Path: %s", err)
	}

	jpath, err = pkg.JPath()
	if err != nil {
		t.Error(err)
	}

	sort.Strings(jpath)
	index = sort.SearchStrings(jpath, examplePath)
	if index == len(jpath) {
		t.Errorf("Unable to find force set path in JSonnet Path")
		t.Errorf("Current %s: %s", JSonnetPath, jpath)
	}
}

func TestArchival(t *testing.T) {
	expected := map[string]string{
		"antho.json":     "a5066f743ea70d60b692ceea51cf8b16a77d1b1c",
		"main.libsonnet": "0b56d40c0630a74abec5398e01c6cd83263feddc",
	}
	received := make(map[string]string)

	buf := new(bytes.Buffer)
	pkg, err := fixturePackage("cache/github.com/seryl/example-0.0.1")
	if err != nil {
		t.Error(err)
	}

	if err := pkg.Build(buf); err != nil {
		t.Error(err)
	}

	tr := tar.NewReader(buf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}

		if err != nil {
			t.Error(err)
		}

		h := sha1.New()
		if _, err := io.Copy(h, tr); err != nil {
			t.Error(err)
		}
		received[hdr.Name] = hex.EncodeToString(h.Sum(nil))
	}

	for k, v := range expected {
		if received[k] != v {
			t.Errorf("Sha1 for `%s` did not match. Expected: `%s` Received: `%s`",
				k, expected[k], received[k])
		}
	}

	if len(expected) != len(received) {
		t.Errorf("Sha1 map does not match. Expected: `%s`, Received: `%s`",
			expected, received)
	}
}

func fixturePackage(name string) (*Package, error) {
	return FromFile(path.Join(fixturesDirectory, name))
}
