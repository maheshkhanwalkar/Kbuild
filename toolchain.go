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
	cflags  string
	ldflags string
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

/*
Compile the given source file
*/
func (tool *Toolchain) Compile(source string) {
	obj := buildObjectPath(source)
	color.Green("[CC] " + getObjectNameFromPath(obj))

	args := append(strings.Fields(tool.cflags), "-c", source, "-o", obj)
	cmd := exec.Command(tool.cc, args...)

	var outBuf, errBuf bytes.Buffer

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	if err != nil {
		fmt.Println(errBuf.String())
		panic("build failed")
	}
}

func (*Toolchain) Link(dir string, objects []string, output string) {
	// TODO
	color.Green("[LD] " + output)
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
