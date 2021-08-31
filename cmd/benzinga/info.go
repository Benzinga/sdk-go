package benzinga

import (
	"fmt"

	"github.com/spf13/cobra"
)

const infoText = "For Support Content licensing@benzinga.com or see https://github.com/Benzinga/sdk-go"

func loadInfoCommands() *cobra.Command {
	command := &cobra.Command{
		Use:   "info",
		Short: "prints into about the benzinga cli",
		Run: func(cmd *cobra.Command, args []string) {
			printInfo()
		},
	}

	return command
}

func printInfo() {
	fmt.Println(infoText)
}
