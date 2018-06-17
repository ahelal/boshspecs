package verify

import (
	"path"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/config"
	"github.com/ahelal/boshspecs/testverifiers"
	log "github.com/sirupsen/logrus"
)

// DownloadAssets this will download assets if needed for the verifier
func DownloadAssets(assets []testverifiers.Asset, testVerifier testverifiers.TestVerifier) error {
	dirAsset, err := config.DirAssets()
	if err != nil {
		return err
	}
	for _, asset := range assets {
		if download, downloadPath := checkAssetNeedDownloading(testVerifier.Name(), dirAsset, asset); download {
			if err := common.DownloadFromURL(asset.DownloadURL, downloadPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkAssetNeedDownloading(verifierName string, assetDir string, asset testverifiers.Asset) (bool, string) {
	downloadPath := path.Join(assetDir, asset.FileName)
	if !common.PathExists(downloadPath) {
		log.Debugf("Asset for %s/%s not found, Download needed", verifierName, asset.FileName)
		return true, downloadPath
	}
	actualShasum, _ := common.Sha2sumFile(downloadPath)
	if actualShasum != asset.Sha256Sum {
		log.Debugf("Asset for %s/%s checksum mismatched. Found '%s' expected '%s', Download needed", verifierName, asset.FileName, asset.Sha256Sum, actualShasum)
		return true, downloadPath
	}
	log.Debugf("Asset for %s/%s found and checksum matched, skipping download", verifierName, asset.FileName)
	return false, ""
}
