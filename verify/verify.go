package verify

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahelal/boshspecs/common"
	"github.com/ahelal/boshspecs/config"
	"github.com/ahelal/boshspecs/testverifiers"
)

//Verify entry point to run all verification
func Verify(mConfig config.MergedConfig, verbose bool, noColor bool) error {
	var instanceGroups []string
	tv, tvParams, err := initialVerifierSteps(mConfig)
	if err != nil {
		return err
	}

	tvParams.TmpDir, err = createTmpDir(mConfig.ConfSpec.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tvParams.TmpDir)

	tvParams.Verbose = verbose
	tvParams.NoColor = noColor
	// Check if we have no filters on instance groups
	if len(mConfig.ConfSpec.Filters.InstanceGroup) == 0 {
		instanceGroups, err = instaceGroups(mConfig)
		if err != nil {
			return err
		}
		for _, instanceGroup := range instanceGroups {
			mConfig.ConfSpec.Filters.InstanceGroup = instanceGroup
			if err := runVerify(mConfig, tv, tvParams); err != nil {
				return err
			}
		}
		return nil
	}
	// instance group defined
	return runVerify(mConfig, tv, tvParams)
}

func runVerify(mConfig config.MergedConfig, tv testverifiers.TestVerifier, tvParams testverifiers.TestVerifierParams) error {
	common.InfoPrint(" targeting instace group '" + mConfig.ConfSpec.Filters.InstanceGroup + "'")
	if mConfig.ConfSpec.LocalExec {
		return verifierLocal(mConfig, tv, tvParams)
	}
	return verifierRemote(mConfig, tv, tvParams)
}

// For now that data is static
func getTargetInfo(spec config.CSpec) (testverifiers.TestVerifierParams, error) {
	platformStatic := testverifiers.TargetPlatform{OS: "linux", Arch: "amd64"}
	return testverifiers.TestVerifierParams{Platform: platformStatic, CSpec: spec}, nil
}

func initialVerifierSteps(mConfig config.MergedConfig) (testverifiers.TestVerifier, testverifiers.TestVerifierParams, error) {
	var tvParams testverifiers.TestVerifierParams
	tv, err := getTestVerifier(mConfig.ConfSpec.Name, mConfig.ConfSpec.SpecType)
	if err != nil {
		return tv, tvParams, err
	}

	mConfig.ConfSpec.Path, err = expandTestPath(mConfig)
	if err != nil {
		return tv, tvParams, err
	}

	tvParams, err = getTargetInfo(mConfig.ConfSpec)
	if err := tv.ValidateConfig(&tvParams); err != nil {
		return tv, tvParams, err
	}

	assets := tv.CheckAssets(&tvParams)
	if assets != nil {
		if err := DownloadAssets(assets, tv); err != nil {
			return tv, tvParams, err
		}
	}
	tvParams.Assets = assets
	tvParams.RemoteTestPath = remoteTestPath + "/test"
	tvParams.RemoteAssetPath = remoteAssetPath

	// Render RunScriptHelperTemplate
	tvParams.RunScriptHelper, err = testverifiers.RenderRunScriptHelperTemplate(&tvParams)
	if err != nil {
		return tv, tvParams, err
	}
	return tv, tvParams, nil
}

//expandTestPath
func expandTestPath(mConfig config.MergedConfig) (string, error) {
	trimmedPath := strings.Trim(mConfig.ConfSpec.Path, " ")

	if len(trimmedPath) == 0 {
		basePath := filepath.Join(config.DirTest, mConfig.ConfSpec.Name)
		return filepath.Abs(basePath)
	}
	if filepath.IsAbs(trimmedPath) {
		return trimmedPath, nil
	}

	return filepath.Abs(mConfig.ConfSpec.Path)
}

func createTmpDir(specName string) (string, error) {
	dirPath, err := ioutil.TempDir(config.DirTMP, specName)
	if err != nil {
		return "", nil
	}
	common.CreateDir(dirPath, "/test")
	return dirPath, nil
}
