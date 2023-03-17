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
func Rake(dir string, config map[string]string, toolchain *Toolchain, action Action) []string {
	var objMap = make(map[string]struct{})
	rake(dir, config, objMap, toolchain, action)

	return toSlice(objMap)
}

func rake(dir string, config map[string]string, objMap map[string]struct{}, toolchain *Toolchain, action Action) {
	kbuild := dir + "/Kbuild"
	denySet := make(map[string]struct{})

	sources := parseKbuild(dir, kbuild, config, denySet)

	for _, source := range sources {
		switch action {
		case Build:
			obj := toolchain.Compile(source)
			objMap[obj] = struct{}{}
		case Clean:
			toolchain.Clean(source)
		}
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		panic("cannot read directory: " + dir)
	}

	for _, entry := range entries {
		if entry.IsDir() && !containsKey(denySet, entry.Name()) {
			child := dir + "/" + entry.Name()
			rake(child, config, objMap, toolchain, action)
		}
	}
}

func parseKbuild(dir string, kbuild string, config map[string]string, denySet map[string]struct{}) []string {
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

		sourceFiles, dirs := categorise(dir, strings.Fields(pieces[1]))
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
			} else {
				/*
					By default, we descend into subdirectories, so if the conditional
					check resolves to false, then explicitly mark the subdirectories
					are denied so rake() will ignore them
				*/
				addToDenySet(dirs, denySet)
			}
		}
	}

	var sources []string

	for _, source := range filesToProcess {
		sources = append(sources, dir+"/"+source)
	}

	return sources
}

func categorise(dir string, input []string) ([]string, []string) {
	var sources []string
	var dirs []string

	for _, elem := range input {
		path := dir + "/" + elem
		info, err := os.Stat(path)

		if err != nil {
			panic("could not stat: " + path)
		}

		if info.IsDir() {
			count := strings.Count(elem, "/")

			if count > 1 || (count > 0 && elem[len(elem)-1] != '/') {
				panic("invalid Kbuild configuration: " + elem + ", only one directory level is allowed")
			}

			// Normalise dir path to not have terminating slash
			if elem[len(elem)-1] == '/' {
				elem = elem[:len(elem)-2]
			}

			dirs = append(dirs, elem)
		} else {
			sources = append(sources, elem)
		}
	}

	return sources, dirs
}

func addToDenySet(slice []string, deny map[string]struct{}) {
	for _, elem := range slice {
		deny[elem] = struct{}{}
	}
}

func containsKey(set map[string]struct{}, key string) bool {
	_, ok := set[key]
	return ok
}

func toSlice(mp map[string]struct{}) []string {
	var res = make([]string, 0, len(mp))

	for key, _ := range mp {
		res = append(res, key)
	}

	return res
}
