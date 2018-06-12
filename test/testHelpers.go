package testhelpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

//RandStringBytes return a random N strings
func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//WriteFile a byte of data into fileName return error if something went wrong
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

//IsDir return true of dir path is a dir and error if did not manage to do a stat call
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

//ReadFile return content of a file as a string
func ReadFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
