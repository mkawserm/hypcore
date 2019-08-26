package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/xcore"
	//"github.com/mkawserm/hypcore/z"
	"github.com/spf13/cobra"
	"os"
)

var hypCoreRootCmd = &cobra.Command{
	Use:   "hypcore",
	Short: "Hyper Core micro service",
	Long:  `Hyper Core is a small reusable golang package to build micro services`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hyper Core",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(aurora.Green(xcore.VERSION))
	},
}

var authorsCmd = &cobra.Command{
	Use:   "authors",
	Short: "Print the authors Hyper Core",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(aurora.Green(xcore.AUTHORS))
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Hyper Core shell",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	hypCoreRootCmd.AddCommand(shellCmd)
	hypCoreRootCmd.AddCommand(authorsCmd)
	hypCoreRootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := hypCoreRootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
