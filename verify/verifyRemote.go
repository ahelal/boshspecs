package verify

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/config"
	"github.com/ahelal/boshspecs/runner"
	"github.com/ahelal/boshspecs/testverifiers"
	"github.com/mholt/archiver"
	log "github.com/sirupsen/logrus"
)

const remoteTestPath = "/var/vcap/boshSpecs/test"
const remoteAssetPath = "/var/vcap/boshSpecs/assets"

func verifierRemote(mConfig config.MergedConfig, tv testverifiers.TestVerifier, tvParams testverifiers.TestVerifierParams) error {
	var stdoutBuf, stderrBuf bytes.Buffer
	log.Debugf("Render and write control script")
	contentOfControlScript, err := renderControlScript(controlScriptType{remoteTestPath, remoteAssetPath, runScriptName, "", tvParams.Assets})
	if err != nil {
		return err
	}

	if _, err = writeControlScript(contentOfControlScript, tvParams.TmpDir); err != nil {
		return err
	}

	log.Debugf("Render and write run script")
	contentOfRunScript, err := tv.GenerateRunScript(&tvParams, "")
	if err != nil {
		return nil
	}
	//TODO refactor pull runscript to initial phase
	runScriptPath := path.Join(tvParams.TmpDir, runScriptName)
	if err = common.WriteFile(runScriptPath, []byte(contentOfRunScript), 0755); err != nil {
		return err
	}

	log.Debugf("Copy tests to tmp dir to prepare for tar gunziping")
	if err := common.CopyToDir(tvParams.CSpec.(config.CSpec).Path, tvParams.TmpDir+"/test"); err != nil {
		return err
	}

	tarGZfileName := path.Base(tvParams.TmpDir) + ".tgz"
	tarGzPath := path.Join(config.DirTMP, tarGZfileName)
	log.Debugf("tgz %s to %s", tvParams.TmpDir, tarGzPath)
	if err := archiver.TarGz.Make(tarGzPath, []string{tvParams.TmpDir}); err != nil {
		return err
	}

	//defer os.RemoveAll(tvParams.TmpDir)
	mConfig.ConfSpec.Filters.Instances = []string{"0"}
	log.Debugf("Remove,create remote dir %s and set permission", remoteTestPath)
	remoteCMD := runner.BoshCMD{Bosh: mConfig.ConfBosh, Deployment: mConfig.ConfigDeployment.Name, InstanceGroup: mConfig.ConfSpec.Filters.InstanceGroup, InstancesIndex: mConfig.ConfSpec.Filters.Instances, Command: remoteInitCmd()}
	if err := verifyboshSSHCommand(remoteCMD, &stdoutBuf, &stderrBuf, false); err != nil {
		return err
	}

	log.Debugf("Uploading tests and control script")
	rscp := runner.BoshCMD{Bosh: mConfig.ConfBosh, Deployment: mConfig.ConfigDeployment.Name, InstanceGroup: mConfig.ConfSpec.Filters.InstanceGroup, InstancesIndex: mConfig.ConfSpec.Filters.Instances, Source: tarGzPath, Dest: remoteTestPath}
	if err := verifyboshSCPCommand(rscp, &stdoutBuf, &stderrBuf); err != nil {
		fmt.Println(stdoutBuf.String(), stderrBuf.String())
		return err
	}

	log.Debugf("untar and execute control script")
	unTarcmd := fmt.Sprintf("tar xf %s/%s -C %s --strip-components=1 && /bin/sh %s/%s", remoteTestPath, tarGZfileName, remoteTestPath, remoteTestPath, controlScriptName)
	remoteCMD.Command = unTarcmd
	err = verifyboshSSHCommand(remoteCMD, &stdoutBuf, &stderrBuf, tvParams.Verbose)
	assetsNeeded, extractErr := extractAssetsDownloadIfAny(&stdoutBuf)
	if err != nil && len(assetsNeeded) == 0 {
		if !tvParams.Verbose {
		}
		return fmt.Errorf("%s. %s", extractErr, err)
	}

	if len(assetsNeeded) > 0 {
		for _, assetNeeded := range assetsNeeded {
			i := getAssetURL(assetNeeded, tvParams.Assets)
			if i == -1 {
				return fmt.Errorf("Unknown asset to download %s", assetNeeded)
			}
			rscp.Source = path.Join(config.DirAssets, tvParams.Assets[i].FileName)
			rscp.Dest = remoteAssetPath
			if err := verifyboshSCPCommand(rscp, &stdoutBuf, &stderrBuf); err != nil {
				return err
			}
		}
		// log.Debugf("Execute control script after assets upload")
		remoteCMD.Command = fmt.Sprintf("/bin/sh %s/%s", remoteTestPath, controlScriptName)
		err := verifyboshSSHCommand(remoteCMD, &stdoutBuf, &stderrBuf, tvParams.Verbose)
		if err != nil && tvParams.Verbose {
			fmt.Println(stdoutBuf.String(), stderrBuf.String())
		}
		return err
	}
	return nil
}

func extractAssetsDownloadIfAny(stderrBuf *bytes.Buffer) ([]string, error) {
	var AssetsToDownload []string
	for _, line := range strings.Split(stderrBuf.String(), "\n") {
		re, err := regexp.Compile(`(>AssetNotFound )(.*)`) // Prepare our regex
		if err != nil {
			return nil, err
		}
		resultSlice := re.FindStringSubmatch(line)
		if len(resultSlice) == 3 {
			AssetsToDownload = append(AssetsToDownload, strings.TrimSpace(resultSlice[2]))
		}
	}
	return AssetsToDownload, nil
}

func getAssetURL(assetDesc string, assets []testverifiers.Asset) int {
	log.Debugf("getAssetURL desc: %s", assetDesc)
	for i, asset := range assets {
		if asset.Description == assetDesc {
			return i
		}
	}
	return -1
}

func uploadTestVerifierAssets(assetsNeeded []string, assets []testverifiers.Asset) error {
	// var stdoutBuf, stderrBuf bytes.Buffer
	for _, assetNeeded := range assetsNeeded {
		i := getAssetURL(assetNeeded, assets)
		if i == -1 {
			return fmt.Errorf("Unkown asset to download %s", assetNeeded)
		}
		// log.Debugf("Uploading asset %s", assets[i]["desc"])
		// common.InfoPrint(fmt.Sprintf("Uploading asset %s", assets[i]["desc"]))
		// if err := verifyboshSCPCommand(lineItem, path.Join(global.DirAssets, assets[i]["file"]), path.Join(remoteAssetPath, assets[i]["file"]), &stdoutBuf, &stderrBuf); err != nil {
		// 	return err
		// }
	}
	return nil
}

/*
runs boot strap remote command
*/
func remoteInitCmd() string {
	msg := `TEST_DIR=%s
			ASSET_DIR=%s
			echo "Removing ${TEST_DIR}"
			sudo rm -rf ${TEST_DIR};
			sudo /bin/mkdir -p {${TEST_DIR},${ASSET_DIR}};
			sudo chown vcap:vcap {${TEST_DIR},${ASSET_DIR}};
			sudo chmod 775 {${TEST_DIR},${ASSET_DIR}}`
	return fmt.Sprintf(msg, remoteTestPath, remoteAssetPath)
}
