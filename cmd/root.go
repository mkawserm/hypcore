package cmd

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/app"
	"github.com/mkawserm/hypcore/xcore"
	"github.com/mkawserm/hypcore/z"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
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

var runServerCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run the main server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var checkConfigFileCmd = &cobra.Command{
	Use:   "check",
	Short: "Check config file",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := ""
		fileNameWithoutExt := ""

		file, err := cmd.Flags().GetString("file")
		if err == nil {
			if file != "" {
				filePath, fileNameWithoutExt = path.Split(file)
				if strings.ToLower(path.Ext(fileNameWithoutExt)) != ".toml" {
					fmt.Println(aurora.BrightRed("Configuration file with .toml extension is allowed"))
					fmt.Println(aurora.BrightRed(path.Ext(fileNameWithoutExt) + " extension is not allowed"))
					return
				}

				fileNameWithoutExt = strings.TrimSuffix(fileNameWithoutExt, path.Ext(fileNameWithoutExt))
			} else {
				filePath = app.ConfigFilePathFirst
				fileNameWithoutExt = app.ConfigFileNameWithoutExt
			}

			v := app.CheckConfigFile(filePath, fileNameWithoutExt)
			if v != nil && app.IsConfigurationOk(v, false) {
				fmt.Println(aurora.BrightGreen("Configuration file Ok"))
				return
			} else {
				fmt.Println(aurora.BrightRed("Invalid configuration file"))
				return
			}
		}
	},
}

var createConfigFileCmd = &cobra.Command{
	Use:   "create",
	Short: "Create config file",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := ""
		fileNameWithoutExt := ""

		var e error
		file, err := cmd.Flags().GetString("file")
		if err == nil {
			if file != "" {
				filePath, fileNameWithoutExt = path.Split(file)
				if strings.ToLower(path.Ext(fileNameWithoutExt)) != ".toml" {
					fmt.Println(aurora.BrightRed("Configuration file with .toml extension is allowed"))
					fmt.Println(aurora.BrightRed(path.Ext(fileNameWithoutExt) + " extension is not allowed"))
					return
				}

				fileNameWithoutExt = strings.TrimSuffix(fileNameWithoutExt, path.Ext(fileNameWithoutExt))
			} else {
				filePath = app.ConfigFilePathFirst
				fileNameWithoutExt = app.ConfigFileNameWithoutExt
			}
			e = app.EnsureDir(filePath)
			if e != nil {
				fmt.Println(aurora.BrightRed(e.Error()))
				return
			}

			v := app.NewConfigFile(filePath, fileNameWithoutExt)
			_ = ioutil.WriteFile(filePath+"/"+fileNameWithoutExt+".toml", []byte(""), 0777)
			e = v.WriteConfig()
			if e != nil {
				fmt.Println(aurora.BrightRed(e.Error()))
				return
			} else {
				fmt.Println(aurora.BrightGreen("Configuration file created"))
				return
			}
		}
	},
}

var configFileCmd = &cobra.Command{
	Use:   "config",
	Short: "create or check configuration file",
	Run: func(cmd *cobra.Command, args []string) {

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
