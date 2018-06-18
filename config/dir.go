package config

import (
	"os"
	"path"

	"github.com/ahelal/boshspecs/common"
)

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
	var homeDir string
	// cUser, err := user.Current()
	// if err != nil {
	// 	return "", err
	// }
	// homeDir := cUser.HomeDir

	// Use directly $HOME for now since cUser.HomeDir checks /etc/passwd
	homeDir = os.Getenv("HOME")
	if len(homeDir) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		homeDir = cwd

	}
	return path.Join(homeDir, dirName), nil
}
