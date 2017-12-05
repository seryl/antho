package antho

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// CacheDirectory is the path used for any package for
// it's linked cache as well as the users global cache.
const CacheDirectory = ".antho"

// GlobalCachePath is the cache path for a users global cache.
func GlobalCachePath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, CacheDirectory), nil
}

/*
Cache is a repository of jsonnet packages and their versions.
There are a number of repositories we want to support:
	* Github
	* Bitbucket
	* Http (Public)
	* Http (Private)
*/
type Cache struct {
	Path string
}

/*
Packages returns a multi-level map for doing lookups on packages
and their versions.

Lookup occurs as [package_name] [version] PACKAGE.
*/
func (c *Cache) Packages() (map[string]map[string]*Package, error) {
	pkgmap := make(map[string]map[string]*Package)

	err := filepath.Walk(
		c.Path, func(fpath string, f os.FileInfo, e error) error {
			if f.IsDir() {
				isPkg, err := IsPackage(fpath)
				if err != nil {
					return err
				}

				if isPkg && !c.isChildCache(fpath) {
					pkg, err := FromFile(fpath)
					if err != nil {
						return err
					}

					// Cleanup package ame to match the download url
					pkgName := strings.TrimPrefix(fpath, c.Path)
					pkgName = strings.TrimSuffix(pkgName,
						fmt.Sprintf("-%s", pkg.Version))
					pkgName = strings.TrimPrefix(pkgName, string(os.PathSeparator))

					// Make map if it's missing
					if pkgmap[pkgName] == nil {
						pkgmap[pkgName] = make(map[string]*Package)
					}
					pkgmap[pkgName][pkg.Version] = pkg
				}
			}
			return nil
		})

	return pkgmap, err
}

func (c *Cache) isChildCache(target string) bool {
	trimmedTarget := strings.TrimPrefix(target, c.Path)
	pathSplit := strings.Split(trimmedTarget, string(os.PathSeparator))
	for _, p := range pathSplit {
		if p == CacheDirectory {
			return true
		}
	}
	return false
}

// LinkedCache is used to symlink items from the global cache
// to your local package jpath.
type LinkedCache struct {
	Path string
}
