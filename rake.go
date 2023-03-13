package main

import (
	"os"
	"strings"
)

type Action int

const (
	Build Action = iota
	Clean
)

/*
Rake together all the source files and perform the requested action
*/
func Rake(dir string, config map[string]string, toolchain *Toolchain, action Action) {
	objMap := make(map[string]string)
	rake(dir, config, objMap, toolchain, action)
}

func rake(dir string, config map[string]string, objMap map[string]string, toolchain *Toolchain, action Action) {
	kbuild := dir + "/Kbuild"
	sources := parseKbuild(dir, kbuild, config)

	for _, source := range sources {
		switch action {
		case Build:
			toolchain.Compile(source)
		case Clean:
			toolchain.Clean(source)
		}
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		panic("cannot read directory: " + dir)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			child := dir + "/" + entry.Name()
			rake(child, config, objMap, toolchain, action)
		}
	}
}

func parseKbuild(dir string, kbuild string, config map[string]string) []string {
	data, err := os.ReadFile(kbuild)

	if err != nil {
		/* ignore on failure -- it's allowed to not have a Kbuild if there's no source files
		   in the given directory */
		return []string{}
	}

	lines := strings.FieldsFunc(string(data), func(c rune) bool {
		return c == '\n'
	})

	var filesToProcess []string

	for _, line := range lines {
		pieces := strings.Split(line, "+=")

		if len(pieces) != 2 {
			panic("invalid Kbuild format: file: " + kbuild + ", line: " + line)
		}

		decl := pieces[0]

		sourceFiles := strings.Fields(pieces[1])
		objType := strings.TrimSpace(strings.Split(decl, "-")[1])

		// Unconditional 'yes' -- always add it in!
		if objType == "y" {
			filesToProcess = append(filesToProcess, sourceFiles...)
		}

		if strings.Contains(objType, "CONFIG_") {
			resolved := config[objType]

			// Conditional 'yes'
			if resolved == "y" {
				filesToProcess = append(filesToProcess, sourceFiles...)
			}
		}
	}

	var sources []string

	for _, source := range filesToProcess {
		sources = append(sources, dir+"/"+source)
	}

	return sources
}
