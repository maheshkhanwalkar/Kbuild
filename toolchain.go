package main

import "github.com/fatih/color"

type Toolchain struct {
	cc      string
	cflags  string
	ldflags string
}

func (*Toolchain) Compile(source string, object string) {
	// TODO
	color.Green("[CC] " + object)
}

func (*Toolchain) Link(objects []string, output string) {
	// TODO
	color.Green("[LD] " + output)
}
