package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"strings"
)

type Toolchain struct {
	cc      string
	cflags  string
	ldflags string
}

func (tool *Toolchain) Compile(dir string, source string, object string) {
	color.Green("[CC] " + object)

	args := append(strings.Fields(tool.cflags), "-c", dir+"/"+source, "-o", dir+"/"+object)
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
