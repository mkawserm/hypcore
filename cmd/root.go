package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/xcore"
	//"github.com/mkawserm/hypcore/z"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hypcore",
	Short: "Hyper Core micro service",
	Long:  `Hyper Core is a small reusable golang package to build micro service`,
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

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Hyper Core shell",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
