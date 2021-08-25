package main

import (
	"github.com/Benzinga/sdk-go/cmd/benzinga"
)

var version = "dev"

func main() {
	benzinga.Run(version)
}
