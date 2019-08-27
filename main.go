package main

import "github.com/mkawserm/hypcore/cmd"
import "github.com/mkawserm/hypcore/xcore"

func main() {
	xcore.AppName = "HypCore"

	cmd.Setup()
	cmd.Execute()
}
