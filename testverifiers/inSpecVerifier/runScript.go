package inspecverifier

const inSpecRunScript = `{{ .TvParams.RunScriptHelper }}

assetPath={{ .TvParams.RemoteAssetPath }}
testSpec="${testDirPath}/*_spec.rb"
assetBin=inspec

inSpecColor="{{.TvParams.NoColor}}"
if [ "${inSpecColor}" = "false" ]; then
    inSpecColor=""
else
   inSpecColor="--no-color"
fi

for FILE in ${testSpec}; do
    FOUND="FOUND"
done

if [ "${FOUND}x" = "x" ]; then
	echo "No _spec.rb files found in ${testDirPath} exiting without error"
	exit 0
fi

FAILED=
for FILE in ${testDirPath}; do
	if ! ${SUDO} ${assetBin} exec "${FILE}" ${inSpecColor}; then
		FAILED=TRUE
	fi
done

if ! [ "${FAILED}x" = "x" ]; then
	exit 1
fi
exit 0
`
