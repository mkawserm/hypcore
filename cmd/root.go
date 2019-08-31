package cmd

import (
	goFlag "flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/xcore"
	"github.com/mkawserm/hypcore/z"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"os"
)

type HyperCoreCMD struct {
	HyperCoreRootCMD *cobra.Command
	VersionCMD       *cobra.Command

	RunServerCMD *cobra.Command
	ServerCMD    *cobra.Command

	CheckConfigFileCMD  *cobra.Command
	CreateConfigFileCMD *cobra.Command
	ConfigFileCMD       *cobra.Command

	AuthorsCMD *cobra.Command
	ShellCMD   *cobra.Command
}

// Load defaults command if not set explicitly
func (hcc *HyperCoreCMD) LoadDefaultsIfNil() {

	if hcc.HyperCoreRootCMD == nil {
		hcc.HyperCoreRootCMD = &cobra.Command{
			Use:   z.ExeName(),
			Short: xcore.AppNameLong + " micro service",
			Long:  xcore.AppDescription,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(cmd.UsageString())
			},
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				// For cobra + glog flags. Available to all sub commands.
				goFlag.Parse()
			},
		}
	}

	if hcc.VersionCMD == nil {
		hcc.VersionCMD = &cobra.Command{
			Use:   "version",
			Short: "Print the version number of Hyper Core",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(aurora.Green(xcore.Version))
			},
		}
	}

	if hcc.ServerCMD == nil {

		hcc.RunServerCMD = &cobra.Command{
			Use:   "run",
			Short: "Run the main server",
			Run:   RunServerCmdRun,
		}

		hcc.RunServerCMD.Flags().String(
			"config",
			"",
			"Absolute path of the configuration file")

		hcc.ServerCMD = &cobra.Command{
			Use:   "server",
			Short: "The main server",
			Long:  "The main HypCore server",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(cmd.UsageString())
			},
		}

		hcc.ServerCMD.AddCommand(hcc.RunServerCMD)
	}

	if hcc.ConfigFileCMD == nil {
		hcc.CheckConfigFileCMD = &cobra.Command{
			Use:   "check",
			Short: "Check config file",
			Run:   CheckConfigFileCmdRun,
		}
		hcc.CheckConfigFileCMD.Flags().String(
			"file",
			"",
			"Absolute file path to check configuration file")

		hcc.CreateConfigFileCMD = &cobra.Command{
			Use:   "create",
			Short: "Create config file",
			Run:   CreateConfigFileCmdRun,
		}
		hcc.CreateConfigFileCMD.Flags().String(
			"file",
			"",
			"Absolute file path to create configuration file")

		hcc.ConfigFileCMD = &cobra.Command{
			Use:   "config",
			Short: "create or check configuration file",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(cmd.UsageString())
			},
		}

		hcc.ConfigFileCMD.AddCommand(hcc.CheckConfigFileCMD)
		hcc.ConfigFileCMD.AddCommand(hcc.CreateConfigFileCMD)
	}

	if hcc.AuthorsCMD == nil {
		hcc.AuthorsCMD = &cobra.Command{
			Use:   "authors",
			Short: "Print the authors of Hyper Core",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(aurora.Green(xcore.Authors))
			},
		}
	}

	if hcc.ShellCMD == nil {
		hcc.ShellCMD = &cobra.Command{
			Use:   "shell",
			Short: "Hyper Core shell",
			Run:   ShellCmdRun,
		}
	}

}

func (hcc *HyperCoreCMD) Setup() {
	hcc.HyperCoreRootCMD.AddCommand(hcc.ShellCMD)
	hcc.HyperCoreRootCMD.AddCommand(hcc.AuthorsCMD)
	hcc.HyperCoreRootCMD.AddCommand(hcc.VersionCMD)
	hcc.HyperCoreRootCMD.AddCommand(hcc.ServerCMD)
	hcc.HyperCoreRootCMD.AddCommand(hcc.ConfigFileCMD)

	flag.CommandLine.AddGoFlagSet(goFlag.CommandLine)
}

func (hcc *HyperCoreCMD) Execute() {
	if err := hcc.HyperCoreRootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
