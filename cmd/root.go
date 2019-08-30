package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/xcore"
	"github.com/mkawserm/hypcore/z"
	"github.com/spf13/cobra"
	"os"
)

var hypCoreRootCmd = &cobra.Command{
	Use:   z.ExeName(),
	Short: xcore.AppNameLong + " micro service",
	Long:  xcore.AppDescription,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hyper Core",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(aurora.Green(xcore.Version))
	},
}

var runServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run the main server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var checkConfigFileCmd = &cobra.Command{
	Use:   "check",
	Short: "Check config file",
	Run:   CheckConfigFileCmd,
}

var createConfigFileCmd = &cobra.Command{
	Use:   "create",
	Short: "Create config file",
	Run:   CreateConfigFileCmd,
}

var configFileCmd = &cobra.Command{
	Use:   "config",
	Short: "create or check configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

var authorsCmd = &cobra.Command{
	Use:   "authors",
	Short: "Print the authors of Hyper Core",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(aurora.Green(xcore.Authors))
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Hyper Core shell",
	Run:   ShellCmd,
}

func Setup() {
	hypCoreRootCmd.AddCommand(shellCmd)
	hypCoreRootCmd.AddCommand(authorsCmd)
	hypCoreRootCmd.AddCommand(versionCmd)
	hypCoreRootCmd.AddCommand(runServerCmd)

	checkConfigFileCmd.Flags().String(
		"file",
		"",
		"Absolute file path to check configuration file")

	createConfigFileCmd.Flags().String(
		"file",
		"",
		"Absolute file path to create configuration file")

	configFileCmd.AddCommand(checkConfigFileCmd)
	configFileCmd.AddCommand(createConfigFileCmd)
	hypCoreRootCmd.AddCommand(configFileCmd)
}

func Execute() {
	if err := hypCoreRootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
