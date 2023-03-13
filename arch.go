package main

const ArchKey = "CONFIG_ARCH"
const SubArchKey = "CONFIG_ARCH_SUB"
const CompilerKey = "CC"
const CompilerFlagsKey = "CFLAGS"
const LinkerFlagsKey = "LDFLAGS"

type Arch struct {
	arch    string
	subArch *string
}

func (arch *Arch) GetArchPath() string {
	return "arch/" + arch.arch
}

func getArch(config map[string]string) *Arch {
	arch := config[ArchKey]
	subArch, ok := config[SubArchKey]

	if !ok {
		return &Arch{arch, nil}
	} else {
		return &Arch{arch, &subArch}
	}
}

func getToolChain(arch *Arch) *Toolchain {
	path := arch.GetArchPath() + "/Kbuild.bootstrap"

	if arch.subArch != nil {
		path += "." + *arch.subArch
	}

	toolchainConfig := readConfig(path)

	return &Toolchain{cc: toolchainConfig[CompilerKey],
		cflags: toolchainConfig[CompilerFlagsKey], ldflags: toolchainConfig[LinkerFlagsKey]}
}
