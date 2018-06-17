package config

import (
	"os"
	"os/user"
	"path"

	"github.com/ahelal/boshspecs/common"
)

//TODO change from const to function tat retrieves abs path

//DirMain main boshspecs dir
func DirMain() (string, error) {
	return getAbsPath(".boshspecs")
}

//DirAssets assets directory
func DirAssets() (string, error) {
	return getAbsPath(".boshspecs/assets")
}

//DirTMP tmp files dir
func DirTMP() (string, error) {
	return getAbsPath(".boshspecs/tmp")
}

//DirTest tmp files dir
func DirTest() (string, error) {
	CWDdir, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	return path.Join(CWDdir, "test"), nil
}

//InitializeDir create directories
func InitializeDir() error {
	dirs := []func() (string, error){DirMain, DirTest, DirTMP, DirAssets}
	for _, dir := range dirs {
		path, err := dir()
		if err != nil {
			return err
		}
		if err := common.CreateDir(path, ""); err != nil {
			return err
		}
	}
	return nil
}

func getAbsPath(dirName string) (string, error) {
	cUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(cUser.HomeDir, dirName), nil
}
