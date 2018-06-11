package testverifiers

//TargetPlatform instance info
type TargetPlatform struct {
	OS           string // linux or windows
	Distribution string // Ubuntu Centos
	Arch         string // x86_64 or x86_32
	Release      string // i.e. 14.04
}

//TestVerifierParams to pass to all function
type TestVerifierParams struct {
	Verbose         bool
	NoColor         bool
	Assets          []Asset
	Platform        TargetPlatform
	TmpDir          string
	CSpec           interface{}
	RunScriptHelper string
	RemoteTestPath  string
	RemoteAssetPath string
}

//TestVerifier is the interface to implement a verifier
type TestVerifier interface {
	Name() string
	ValidateConfig(*TestVerifierParams) error
	CheckAssets(*TestVerifierParams) []Asset
	GenerateRunScript(*TestVerifierParams, string) (string, error)
}

// Asset describes an asset
type Asset struct {
	Description   string
	DownloadURL   string
	Sha256Sum     string
	FileName      string
	FileMode      string
	InstallCMD    string
	InstallPath   string
	InstallVerify string
	InstallSudo   bool
}
