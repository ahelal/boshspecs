package verify

import (
	"fmt"
	"strings"

	"github.com/ahelal/boshspecs/testverifiers"
	gossVerifier "github.com/ahelal/boshspecs/testverifiers/gossverifier"
	inSpecVerifier "github.com/ahelal/boshspecs/testverifiers/inspecverifier"
	shellVerifier "github.com/ahelal/boshspecs/testverifiers/shellverifier"
)

type assetsType []map[string]string

func getTestVerifier(verifierName string, verifierType string) (testverifiers.TestVerifier, error) {
	switch strings.ToLower(verifierType) {
	case "shell":
		return shellVerifier.ShellTestVerifier{}, nil
	case "goss":
		return gossVerifier.GossTestVerifier{}, nil
	case "inspec":
		return inSpecVerifier.InSpecTestVerifier{}, nil
	}
	return nil, fmt.Errorf("Unknown test verifier '%s' defined with '%s'", verifierType, verifierName)
}
