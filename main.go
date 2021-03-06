package main

import "github.com/mkawserm/hypcore/cmd"
import "github.com/mkawserm/hypcore/package/xcore"

func main() {
	xcore.AppName = "HypCore"

	hcc := new(cmd.HyperCoreCMD)
	hcc.LoadDefaultsIfNil()
	hcc.Setup()
	hcc.Execute()
}
