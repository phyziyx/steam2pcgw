package main

import (
	"fmt"
	"os"
	"steam2pcgw/src/pkg"
)

func main() {
	if err := pkg.Run(os.Args); err != nil {
		fmt.Println("An error occurred:", err)
	}
}
