package config

import (
	"path/filepath"

	"github.com/ahelal/boshspecs/common"
)

//TODO change from const to function tat retrieves abs path

//DirMain main boshspecs dir
const DirMain = ".boshSpecs"

//DirAssets assets directory
const DirAssets = ".boshSpecs/assets"

//DirTMP tmp files dir
const DirTMP = ".boshSpecs/tmp"

//DirTest test path dir
const DirTest = "test"

//InitializeDir create directories
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
