package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/seryl/antho"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func jpathCmd(pkgPathList []string) int {
	var pkgPath string

	wd, err := os.Getwd()
	if err != nil {
		Logger.Error(err)
		return 1
	}

	if len(pkgPathList) < 1 {
		pkgPath = wd
	}

	if len(pkgPathList) > 1 {
		Logger.Error("Only one package path may be supplied")
		return 1
	}

	if len(pkgPathList) == 1 {
		if !filepath.IsAbs(pkgPath) {
			pkgPath = path.Join(pkgPath, pkgPathList[0])
		} else {
			pkgPath = pkgPathList[0]
		}
	}

	pkg, err := antho.FromFile(pkgPath)
	if err != nil {
		Logger.WithFields(log.Fields{
			"path":  pkgPath,
			"error": err,
		}).Error("error reading package")
		return 1
	}

	jpaths, err := pkg.JPath()
	if err != nil {
		Logger.WithFields(log.Fields{
			"path":  pkgPath,
			"error": err,
		}).Error("Unable to retrieve jpaths")
		return 1
	}

	fmt.Println(strings.Join(jpaths, ";"))
	return 0
}

// CmdJPath will print the jsonnet path for a given package.
// Note: This defaults to checking the current directory.
var CmdJPath = &cobra.Command{
	Use:   "jpath PKG",
	Short: "prints the Jsonnet Path for a package",
	Long:  `prints the Jsonnet Path for a package`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(jpathCmd(args))
	},
}
