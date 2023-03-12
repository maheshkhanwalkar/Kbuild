package main

type Toolchain struct {
	cc      string
	cflags  string
	ldflags string
}

func (*Toolchain) Compile(source string, object string) {
	// TODO
}

func (*Toolchain) Link(objects []string, output string) {
	// TODO
}
