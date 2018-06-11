package config

import (
	"path/filepath"

	"github.com/ahelal/boshspecs/common"
)

const DirMain = ".boshSpecs"
const DirAssets = ".boshSpecs/assets"
const DirTMP = ".boshSpecs/tmp"
const DirTest = "test"

func InitializeDir() error {
	for _, path := range []string{DirMain, DirTest, DirTMP, DirAssets} {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if err := common.CreateDir(absPath, ""); err != nil {
			return err
		}
	}
	return nil
}
