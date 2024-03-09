package pkg

import (
	"fmt"
	"steam2pcgw/src/api"
)

// Version command
func VersionCommand(a App, parameters []string) error {
	a.InfoPrint()
	return nil
}

// Article generation command
func GenerateArticleCommand(a App, parameters []string) error {
	if len(parameters) < 1 {
		return fmt.Errorf("please provide a Steam app ID")
	}

	return api.GenerateArticle(parameters[0])
}

// Game cover command
func GenerateCoverCommand(a App, parameters []string) error {
	if len(parameters) < 1 {
		return fmt.Errorf("please provide a Steam app ID")
	}

	return api.GenerateCover(parameters[0])
}
