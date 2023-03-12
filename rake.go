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
	parseKbuild(dir, kbuild, config, objMap, toolchain)

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

func parseKbuild(dir string, kbuild string, config map[string]string, objMap map[string]string, toolchain *Toolchain) {
	data, err := os.ReadFile(kbuild)

	if err != nil {
		/* ignore on failure -- it's allowed to not have a Kbuild if there's no source files
		   in the given directory */
		return
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

	for _, source := range filesToProcess {
		output := []rune(source)
		output[len(output)-1] = 'o'

		object := string(output)
		toolchain.Compile(dir, source, object)

		objMap[source] = string(output)
	}
}
