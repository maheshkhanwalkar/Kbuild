package main

func main() {
	println("[kbuild]")

	config := readConfig(".config")
	arch := getArch(config)
	toolchain := getToolChain(arch, config)

	objs := Rake(arch.GetArchPath(), config, toolchain, Build)
	toolchain.Link(".", objs, "vminix")
}
