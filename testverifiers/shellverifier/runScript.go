package shellverifier

const shellRunScriptTemplate = `{{ .TvParams.RunScriptHelper }}
testShellGlob="${testDirPath}/*.sh"
FAILED=
for FILE in ${testShellGlob}; do
	${SUDO} chmod +x "${FILE}"
	if ! ${SUDO} ${FILE}; then
		FAILED=TRUE
	fi
done

if [ "${FAILED}x" = "x" ]; then
	exitMSG 0 ""
else
	exitMSG 1 ""
fi
exit 0
`
