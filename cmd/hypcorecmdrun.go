package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

const hypCoreVersion = "1.0.0"
const hypCoreAuthors = "Md Kawser Munshi"

func HypCoreCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println(aurora.Green("HypCore Version: " + hypCoreVersion))
	fmt.Println(aurora.Green("HypCore Authors: " + hypCoreAuthors))
}
