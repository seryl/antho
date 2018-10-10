package cmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/seryl/antho/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func testCmd(flags *pflag.FlagSet, pkgPathList []string) int {
	if len(pkgPathList) < 1 {
		Logger.Error("no package paths were supplied")
		return 1
	}

	wd, err := os.Getwd()
	if err != nil {
		Logger.Error(err)
		return 1
	}

	outputDir := wd
	if flags != nil {
		outputTargetDir, err := flags.GetString("output")
		if err != nil {
			Logger.Error(err)
			return 1
		}

		if outputTargetDir != "" {
			if filepath.IsAbs(outputTargetDir) {
				outputDir = outputTargetDir
			} else {
				outputDir = path.Join(wd, outputTargetDir)
			}
		}
	}

	for _, p := range pkgPathList {
		pkgPath := path.Join(wd, p)
		pkg, err := antho.FromFile(pkgPath)

		Logger.WithFields(log.Fields{
			"package": pkg.Name,
			"path":    p,
		}).Info("reading jsonnet package info")

		if err != nil {
			Logger.WithFields(log.Fields{
				"path":  p,
				"error": err,
			}).Error("error reading package")
			return 1
		}

		Logger.WithFields(log.Fields{
			"package": pkg.Name,
			"path":    p,
		}).Info("writing jsonnet tarball")

		err = pkg.WriteTarball(outputDir)
		if err != nil {
			Logger.WithFields(log.Fields{
				"package": pkg.Name,
				"path":    p,
				"error":   err,
			}).Error("error tarballing package")
			return 1
		}

		Logger.WithFields(log.Fields{
			"package": pkg.Name,
			"tarball": pkg.TarballName(),
		}).Info("packaged successfully")
	}

	return 0
}

// CmdTest will evaluate a .
var CmdTest = &cobra.Command{
	Use:   "test PKG",
	Short: "tests a package against it's output yaml",
	Long:  `tests a package against it's output yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(packageTest(cmd.Flags(), args))
	},
}
