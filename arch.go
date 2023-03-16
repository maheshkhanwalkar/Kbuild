package main

import "strings"

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

func getToolChain(arch *Arch, config map[string]string) *Toolchain {
	path := arch.GetArchPath() + "/Kbuild.bootstrap"

	if arch.subArch != nil {
		path += "." + *arch.subArch
	}

	toolchainConfig := readConfig(path)
	macroFlags := buildMacroFlags(config)

	cflags := append(strings.Fields(toolchainConfig[CompilerFlagsKey]), macroFlags...)

	return &Toolchain{cc: toolchainConfig[CompilerKey],
		cflags: cflags, ldflags: strings.Fields(toolchainConfig[LinkerFlagsKey])}
}

func buildMacroFlags(config map[string]string) []string {
	var res []string

	for key, value := range config {
		if value == "n" {
			continue
		}

		if value == "y" {
			res = append(res, "-D"+key)
		} else {
			res = append(res, "-D"+key+"="+value)
		}
	}

	return res
}
