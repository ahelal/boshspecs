package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

//DownloadFromURL download a file from a url
func DownloadFromURL(url string, downloadPath string) error {
	log.Debugf("Downloading %s to %s", url, downloadPath)
	output, err := os.Create(downloadPath)
	if err != nil {
		log.Warnf("Error while creating '%s' %s", downloadPath, err)
		return fmt.Errorf("Error while creating '%s' %s", downloadPath, err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Warnf("Error while downloading  '%s' %s", url, err)
		return fmt.Errorf("Error while downloading  '%s' %s", url, err)
	}
	if response.StatusCode != 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			body = []byte(err.Error())
		}
		return fmt.Errorf("got a non 200. %s", body)
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Warnf("Error while downloading  '%s' %s", url, err)
		return fmt.Errorf("Error while downloading  '%s' %s", url, err)
	}
	return nil
}
