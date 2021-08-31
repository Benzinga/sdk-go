package benzinga

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/Benzinga/sdk-go/internal/cli"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var config cli.Config

func Run(version string) {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		fmt.Println(info)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("logger setup failed: ", err)
	}

	zap.ReplaceGlobals(logger)

	defer logger.Sync()

	rootCmd := &cobra.Command{
		Use:     "benzinga",
		Short:   "a cli for Benzinga services",
		Long:    `benzinga is a CLI to interact with Benzinga services.`,
		Version: version,
	}

	rootCmd.PersistentFlags().BoolVar(&config.Interactive, "no-interactive", false, "disable interactive prompts")
	rootCmd.PersistentFlags().BoolVar(&config.Debug, "debug", false, "enable debug logging")

	newsCommands := loadNewsCommands()
	infoCommands := loadInfoCommands()

	rootCmd.AddCommand(newsCommands, infoCommands)

	if err := rootCmd.Execute(); err != nil {
		zap.L().Error("execution error", zap.Error(err))
	}
}
