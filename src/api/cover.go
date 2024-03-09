package api

import (
	"fmt"
	"strconv"
)

func GenerateCover(appId string) error {
	if _, err := strconv.Atoi(appId); err != nil {
		return fmt.Errorf("an invalid app ID was provided")
	}

	// https: //cdn.cloudflare.steamstatic.com/steam/apps/990080/library_600x900_2x.jpg

	return nil
}
