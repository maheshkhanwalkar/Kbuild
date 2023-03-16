package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

type Toolchain struct {
	cc      string
	cflags  []string
	ldflags []string
}

func buildObjectPath(source string) string {
	chars := []rune(source)
	chars[len(chars)-1] = 'o'
	return string(chars)
}

func getObjectNameFromPath(path string) string {
	pos := strings.LastIndex(path, "/") + 1
	return path[pos:]
}

func execCmd(name string, args []string) {
	cmd := exec.Command(name, args...)

	var outBuf, errBuf bytes.Buffer

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	if err != nil {
		fmt.Println(errBuf.String())
		panic("build failed")
	}
}

/*
Compile the given source file
*/
func (tool *Toolchain) Compile(source string) string {
	obj := buildObjectPath(source)
	color.Green("[CC] " + getObjectNameFromPath(obj))

	args := append(tool.cflags, "-c", source, "-o", obj)
	execCmd(tool.cc, args)

	return obj
}

/*
Link the object files into one output file
*/
func (tool *Toolchain) Link(dir string, objects []string, output string) {
	color.Green("[LD] " + output)

	obj := dir + "/" + output

	args := append(tool.ldflags, objects...)
	args = append(args, "-o", obj)

	execCmd(tool.cc, args)
}

/*
Clean the associated object file for this source file
*/
func (tool *Toolchain) Clean(source string) {
	obj := buildObjectPath(source)
	err := os.Remove(obj)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warn: could not remove %s", obj)
	}
}
