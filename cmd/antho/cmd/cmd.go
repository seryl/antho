package cmd

import (
	"fmt"
	"os"

	"github.com/seryl/antho"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Logger is the logger to use for all commands.
var Logger *log.Logger

// Config is the command viper configuration.
var Config *viper.Viper

// Execute runs the cli app for antho.
func Execute() {
	Config = viper.New()
	Config.SetEnvPrefix("antho")
	Config.AutomaticEnv()

	Logger = InitializeLogger(Config)

	// Configuration Defaults
	Config.SetDefault("debug", false)
	Config.SetDefault("formatter", "text")

	// CLI Commands
	rootCmd := &cobra.Command{Use: "antho"}

	var cmdHelp = &cobra.Command{
		Use:   "help [command]",
		Short: "show help",
		Long: `help provides help for any command in the application.
				simply type ` + rootCmd.Name() + ` help [path to command] for full details.`,
		Run: rootCmd.HelpFunc(),
	}
	rootCmd.SetHelpCommand(cmdHelp)

	// Flags List
	rootCmd.Flags().BoolP("debug", "d", false, "enable debug mode")
	Config.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))

	rootCmd.Flags().StringP("formatter", "f", "text", "the logging formatter to use [text|json]")
	Config.BindPFlag("formatter", rootCmd.Flags().Lookup("formatter"))

	// Output for CmdPackage
	wd, err := os.Getwd()
	if err != nil {
		Logger.Error(err)
		return
	}
	CmdPackage.Flags().StringP("output", "o", wd, "outputs tarballs to the given path")

	// Command List
	rootCmd.AddCommand(CmdPackage)
	rootCmd.AddCommand(CmdJPath)
	rootCmd.AddCommand(CmdVersion)

	if err := rootCmd.Execute(); err != nil {
		Logger.Error(err)
	}
}

// CmdVersion is the version command.
var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "print the version",
	Long:  `print the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(antho.Version)
	},
}
