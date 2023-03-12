package main

func main() {
	println("kbuild")

	config := readConfig(".config")
	arch := getArch(config)
	_ = getToolChain(arch)
}
