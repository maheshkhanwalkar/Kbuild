package main

func main() {
	println("[kbuild]")

	config := readConfig(".config")
	arch := getArch(config)
	toolchain := getToolChain(arch)

	Rake("arch/"+arch.arch, "arch.a", config, toolchain)
}
