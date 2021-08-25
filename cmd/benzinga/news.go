package benzinga

import (
	"github.com/Benzinga/sdk-go/internal/news/rest/export"
	"github.com/spf13/cobra"
)

func loadNewsCommands() *cobra.Command {
	command := &cobra.Command{
		Use: "news",
	}

	// Subcommands

	exportConfig := export.NewConfig(&config)

	exp := &cobra.Command{
		Use:   "export",
		Short: "start export from News API",
		Run: func(cmd *cobra.Command, args []string) {
			export.Start(exportConfig)
		},
	}

	exp.LocalNonPersistentFlags().StringVarP(&exportConfig.OutputDirectory, "dir", "d", "", "writeable directory to place export files")

	command.AddCommand(exp)

	return command
}
