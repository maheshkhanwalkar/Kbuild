package main

import "os"

func main() {
	println("[kbuild]")

	config := readConfig(".config")
	arch := getArch(config)
	toolchain := getToolChain(arch, config)

	var action Action

	if len(os.Args) == 1 {
		action = Build
	} else if os.Args[1] == "clean" {
		action = Clean
	}

	objs := Rake(arch.GetArchPath(), config, toolchain, action)

	if action == Build {
		toolchain.Link(".", objs, "vminix")
	}
}
