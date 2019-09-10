package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/app"
	"github.com/spf13/cobra"
)

func CreateSuperUserCmdRun(cmd *cobra.Command, args []string) {
	userName, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	if userName == "" || password == "" {
		fmt.Println(aurora.BrightRed("username and password cannot be empty"))
		return
	}

	v := Common(cmd, args, false)

	if v != nil {
		app.AddSuperUser(v, userName, password)
	}
}
