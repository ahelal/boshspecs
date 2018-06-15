package config

import (
	"path/filepath"

	"github.com/ahelal/boshspecs/common"
)

//TODO change from const to function tat retrieves abs path

//DirMain main boshspecs dir
const DirMain = ".boshspecs"

//DirAssets assets directory
const DirAssets = ".boshspecs/assets"

//DirTMP tmp files dir
const DirTMP = ".boshspecs/tmp"

//DirTest test path dir
const DirTest = "test"

//InitializeDir create directories
func InitializeDir(basePath string) error {
	for _, path := range []string{DirMain, DirTest, DirTMP, DirAssets} {
		dirToCreate := filepath.Join(basePath, path)
		if err := common.CreateDir(dirToCreate, ""); err != nil {
			return err
		}
	}
	return nil
}
