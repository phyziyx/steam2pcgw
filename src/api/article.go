package api

import (
	"fmt"
	"strconv"
)

func GenerateArticle(appId string) error {
	if _, err := strconv.Atoi(appId); err != nil {
		return fmt.Errorf("an invalid app ID was provided")
	}

	return nil
}
