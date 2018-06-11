package testhelpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func WriteFile(fileName string, data []byte) error {
	err := ioutil.WriteFile(fileName, data, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("Failed to write to file '%s', %s", fileName, err.Error())
	}
	return nil
}

//WriteTmpFilecontent write content to a temp path
func WriteTmpFilecontent(content []byte) string {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(tmpfile.Name(), 0777)
	if err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}

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

func ReadFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
