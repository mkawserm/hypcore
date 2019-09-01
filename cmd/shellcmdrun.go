package cmd

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/package/xcore"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func ShellCmdRun(cmd *cobra.Command, args []string) {
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
}
