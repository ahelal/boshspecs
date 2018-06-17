package verify

import (
	"fmt"
	"strings"

	"github.com/ahelal/boshspecs/testverifiers"
	gossverifier "github.com/ahelal/boshspecs/testverifiers/gossverifier"
	inspecverifier "github.com/ahelal/boshspecs/testverifiers/inspecverifier"
	shellverifier "github.com/ahelal/boshspecs/testverifiers/shellverifier"
)

type assetsType []map[string]string

func getTestVerifier(verifierName string, verifierType string) (testverifiers.TestVerifier, error) {
	switch strings.ToLower(verifierType) {
	case "shell":
		return shellverifier.ShellTestVerifier{}, nil
	case "goss":
		return gossverifier.GossTestVerifier{}, nil
	case "inspec":
		return inspecverifier.InSpecTestVerifier{}, nil
	}
	return nil, fmt.Errorf("Unknown test verifier '%s' defined with '%s'", verifierType, verifierName)
}
