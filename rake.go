package main

import (
	"os"
	"strings"
)

/*
	Rake together all the objects to be built
*/
func Rake(dir string, archive string, config map[string]string, toolchain *Toolchain) {
	objMap := make(map[string]string)
	rake(dir, config, objMap, toolchain)
}

func rake(dir string, config map[string]string, objMap map[string]string, toolchain *Toolchain) {
	kbuild := dir + "/Kbuild"
	parseKbuild(kbuild, config, objMap, toolchain)

	entries, err := os.ReadDir(dir)

	if err != nil {
		panic("cannot read directory: " + dir)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			child := dir + "/" + entry.Name()
			rake(child, config, objMap, toolchain)
		}
	}
}

func parseKbuild(kbuild string, config map[string]string, objMap map[string]string, toolchain *Toolchain) {
	data, err := os.ReadFile(kbuild)

	if err != nil {
		panic("could not read Kbuild, failing: " + err.Error())
	}

	lines := strings.Split(string(data), "\n")
	var filesToProcess []string

	for _, line := range lines {
		pieces := strings.Split(line, "+=")

		if len(pieces) != 2 {
			panic("invalid Kbuild format: file: " + kbuild + ", line: " + line)
		}

		decl := pieces[0]
		sourceFiles := strings.Split(pieces[0], " ")

		objType := strings.Split(decl, "-")[1]

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

	for _, source := range filesToProcess {
		output := []rune(source)
		output[len(output)-1] = 'o'

		object := string(output)
		toolchain.Compile(source, object)

		objMap[source] = string(output)
	}
}
