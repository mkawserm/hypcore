package cmd

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/xcore"
	"github.com/mkawserm/hypcore/z"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

var authorsCmd = &cobra.Command{
	Use:   "authors",
	Short: "Print the authors Hyper Core",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(aurora.Green(xcore.Authors))
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Hyper Core shell",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		input_counter := 0
		for {
			input_counter++
			fmt.Printf("%s%s%d%s%s ",
				aurora.Bold(aurora.Green(xcore.AppName)),
				aurora.Bold(aurora.Green("[")),
				aurora.Bold(aurora.Red(input_counter)),
				aurora.Bold(aurora.Green("]")),
				"$",
			)
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			cmdString = strings.TrimSuffix(cmdString, "\n")

			switch cmdString {
			case "clear":
				fmt.Print("\x1b[H\x1b[2J")
			case "version":
				fmt.Println(aurora.Green(xcore.Version))
			case "authors":
				fmt.Println(aurora.Green(xcore.Authors))
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
