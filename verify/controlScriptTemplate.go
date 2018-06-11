package verify

const controlScriptName = "controlScript.sh"
const controlScriptTemplate = `#!/bin/sh
set -eu

shasum="/usr/bin/shasum -a 256"

#* basePath
export testPath={{.TestPath}}
export assetPath={{.AssetPath}}
export runScriptPath=${testPath}/{{.RunScript}}

missingBinary(){
    # 1 fileDesc, 2 reason
    echo "${2}"
    echo ">AssetNotFound ${1}"
}

replace(){
    STR="${1}"
    SUB="${2}"
    REP="${3}"
    echo "${STR}" | sed "s:${SUB}:${REP}:g"
}

checkFile(){
  filePath="${1}"
  fileDesc="${2}"
  fileSum="${3}"
  installSudo="${4}"
  installCMD="${5}"
  installPath="${6}"
  InstallVerify="${7}"
  
  # Check if file exists
  if ! [ -f "${filePath}" ]; then
    missingBinary "${fileDesc}" "${filePath} Not found."
    return 1
  fi
    
  resultShasum=$(${shasum} "${filePath}" | awk '{print $1}' 2> /dev/null)
  if [ "${resultShasum}" != "${fileSum}" ]; then
    missingBinary "${fileDesc}" "${filePath} shasum mismatch."
    return 1
  fi

  if [ "${installCMD}x" = "x" ]; then
    echo "> making ${filePath} executable"
    sudo chmod +x "${filePath}"
  else
	InstallVerify=$(replace "${InstallVerify}" "<filePath>" "${filePath}")
    if ${installSudo} ${InstallVerify} > /dev/null; then
      echo "> ${filePath} installed "
    else
      echo "> Installing ${filePath}"
      installCMD=$(replace "${installCMD}" "<filePath>" "${filePath}")
      if ! ${installSudo} ${installCMD}; then
        echo "> installation failed"
        exit 1
      fi
    fi
  fi
  return 0
}

randomString(){
  echo $(od -vAn -N4 -tu4 < /dev/urandom | awk '{print $1}' | base64)
}


# Main
fileCopyRequired="/tmp/$(randomString)-copyRequired"
{{if .Assets}} 
{{	range $i, $asset := .Assets -}}

export check_path="${assetPath}/{{$asset.FileName}}"
export check_desc="{{$asset.Description}}"
export check_sum="{{$asset.Sha256Sum}}"
# Install options if any
export install_cmd="{{$asset.InstallCMD}}"
export install_path="{{$asset.InstallPath}}"
export install_verify="{{if $asset.InstallVerify }}{{$asset.InstallVerify}}{{else}}/bin/true{{end}}"
export install_sudo="{{if $asset.InstallSudo }}sudo{{end}}"

# Call to check script exports
if ! checkFile "${check_path}" "${check_desc}" "${check_sum}" "${install_sudo}" "${install_cmd}" "${install_path}" "${install_verify}" ; then
  touch "${fileCopyRequired}"
fi

{{ end }}
if [ -f "${fileCopyRequired}" ]; then
    rm -f "${fileCopyRequired}"
    echo ">AbortExec missing assets"
    exit 1
fi
rm -f "${fileCopyRequired}"

{{else}}
# No assets
{{end}}

echo ">RuningTests"
chmod +x ${runScriptPath}
${runScriptPath}

echo "Execution went through no errors"
exit 0
`
