package main

import (
	"os"
	"strings"
)

/*
	Read the kernel config and build a key-value mapping
*/
func readConfig(file string) map[string]string {
	data, err := os.ReadFile(file)

	if err != nil {
		panic("could not read config!")
	}

	lines := strings.Split(string(data), "\n")
	result := make(map[string]string)

	for _, line := range lines {
		kv := strings.Split(line, "=")

		if len(kv) != 2 {
			panic("error: invalid config format: " + line)
		}

		key := kv[0]
		value := kv[1]

		result[key] = value
	}

	return result
}
