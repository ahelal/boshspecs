package common

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	otiai "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

//PathExists check if file or path does not exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// ReadFile return content of a file as bytes
func ReadFile(fileName string) ([]byte, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// WriteFile write an array of bytes to fileName
func WriteFile(fileName string, data []byte, fileMode os.FileMode) error {
	err := ioutil.WriteFile(fileName, data, fileMode)
	if err != nil {
		return fmt.Errorf("Failed to write to file '%s', %s", fileName, err.Error())
	}
	return nil
}

// Sha2sumFile return a sha2sum of a file
func Sha2sumFile(filePath string) (string, error) {
	// Open file for reading
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return "", err
	}

	// Create new hasher, which is a writer interface
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	sum := hasher.Sum(nil)
	sumHex := fmt.Sprintf("%x", sum)
	return sumHex, nil
}

// GetCWD return curret working directory
func GetCWD() string {
	dir, _ := os.Getwd()
	return dir
}

// CreateDir Create directories if needed
func CreateDir(baseDir string, requestedDir string) error {
	dirPath := path.Join(baseDir, requestedDir)
	if PathExists(dirPath) {
		return nil
	}
	if err := os.Mkdir(dirPath, 0750); err != nil {
		return err
	}
	return nil
}

//RmDir Removes a dir content and a dir
func RmDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//IsDir return true if dirPath is a directory
func IsDir(dirPath string) (bool, error) {
	fi, err := os.Stat(dirPath)
	if err != nil {
		return false, err
	}
	mode := fi.Mode()
	if mode.IsDir() {
		return true, nil
	}
	return false, nil
}

//CopyToDir src path to dest path
func CopyToDir(src string, destDir string) error {
	log.Debugf("Copy src '%s' to '%s'", src, destDir)
	if err := otiai.Copy(src, destDir); err != nil {
		return err
	}

	return nil
}
