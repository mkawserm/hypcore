package main

import "github.com/mkawserm/hypcore/cmd"
import "github.com/mkawserm/hypcore/xcore"

func main() {
	xcore.AppName = "HypCore"

	hcc := new(cmd.HyperCoreCMD)
	hcc.LoadDefaults()
	hcc.Setup()
	hcc.Execute()
}
