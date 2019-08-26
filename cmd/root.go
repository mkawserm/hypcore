package cmd

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/xcore"
	"strings"

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
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print(aurora.Bold(aurora.Green("HypCore$ ")))
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			cmdString = strings.TrimSuffix(cmdString, "\n")

			switch cmdString {
			case "clear":
				fmt.Print("\x1b[H\x1b[2J")
			case "version":
				fmt.Println(aurora.Green(xcore.VERSION))
			case "exit":
				os.Exit(1)
			}
		}

	},
}

func Setup() {
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
