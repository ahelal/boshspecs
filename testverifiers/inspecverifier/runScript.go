package inspecverifier

const inSpecRunScript = `{{ .TvParams.RunScriptHelper }}

FOUND=""

# If running locally copy test folders before doing anything
{{ if .TvParams.CSpec.LocalExec }}
tmpDir=$(mktemp -d)
cp -r ${testDirPath}/. ${tmpDir}
testDirPath="${tmpDir}"
{{ end }}

specGlob="*_spec.rb"
assetBin=inspec

INSPEC_YAML="
name: boshspec-default
title: default profile
version: 0.1.0
depends:
  - name: inspec-bosh
    url: https://github.com/ahelal/inspec-bosh/archive/master.tar.gz
"

inSpecColor="{{.TvParams.NoColor}}"
if [ "${inSpecColor}" = "false" ]; then
    inSpecColor=""
else
   inSpecColor="--no-color"
fi

# Create the controls dir
mkdir -p "${testDirPath}/controls"

# Find might be problamtic have a look at https://github.com/koalaman/shellcheck/wiki/SC2044
for f in $(find ${testDirPath} -maxdepth 1 -name ${specGlob}); do
	FOUND="FOUND"
	echo "Moving ${f}" to "${testDirPath}/controls" 
	mv "${f}" "${testDirPath}/controls"
done

if [ -d "${testDirPath}/controls" ]; then
	for f in $(find ${testDirPath}/controls -name ${specGlob}); do
		FOUND="$f"
	done
fi

# Manage inspec.yml
if [ -f "${testDirPath}/inspec.yml" ]; then
	FOUND="FOUND"
else
	echo "creating default ${testDirPath}/inspec.yml"
	echo "${INSPEC_YAML}" > "${testDirPath}/inspec.yml"
fi

if [ "${FOUND}x" = "x" ]; then
	echo "No ${specGlob} files found in ${testDirPath} exiting without error"
	exit 0
else 
	echo "Found rb files will run tests in ${testDirPath}"
fi

${SUDO} ${assetBin} exec "${testDirPath}" ${inSpecColor}

exit 0
`
