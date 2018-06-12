package gossverifier

const gossRunScript = `{{ .TvParams.RunScriptHelper }}

assetPath={{ .TvParams.RemoteAssetPath }}
testYML="${testDirPath}/*.yml"

gossColor="{{.TvParams.NoColor}}"
if [ "${gossColor}" = "false" ]; then
    gossColor=""
else
   gossColor="--no-color"
fi

for FILE in ${testYML}; do
    FOUND="FOUND"
done

if [ "${FOUND}x" = "x" ]; then
    echo "No yaml files found in ${testYML} exiting without error"
    exit 0
fi

sudo cat ${testYML} | sudo ${assetPath}/goss-linux-amd64 -g - v ${gossColor}
exit 0
`
