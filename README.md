# Kbuild
Inix Kernel Build System

This project defines the source for the Inix Kernel build system. It is heavily inspired (in principle)
by the Linux kernel build system, which it shares a name with, but takes a drastically different implementation
approach.

There are three main components of the build system:

1. **.config**


    The .config file in the kernel root contains a key-value pair configuration file which specifies which
    options should be selected -- i.e. how the kernel should be configured.

    Within the .config, one configuration is absolutely mandatory:

        CONFIG_ARCH=<arch-here>
        CONFIG_ARCH_SUB=<sub-arch-here> (optional)

    The architecture configuration informs Kbuild which arch-specific directory to descend into to build those
    components - but more crucially - what compiler and flags to use for that architecture!


2. **arch/{arch}/Kbuild.bootstrap.{subarch}**


    The Kbuild bootstrap file specifies the compiler, CFLAGS, and LDFLAGS required for the architecture. This
    will be applied globally to all files that will be compiled by the build system.

3. **Kbuild**


    The Kbuild file in a kernel source directory specifies what source files to compile. It can have two kinds
    of statements:

        src-y += kernel.c init.c ...
        src-${CONFIG_XYZ} += xyc.c ...

    The first statement is unconditional -- the 'y' indicates that those source files must always be included,
    regardless of configuration. The second statement is conditional -- if CONFIG_XYZ=y in the .config file,
    only then will the source files be included.


Rather than leveraging Makefiles or CMake, Kbuild is implemented from scratch in
Go. This leads to a clean and elegant implementation rather than trying to hack
things to work in existing build systems.

For the inix kernel, this build system is sufficient -- we have no need for complex
build logic that would necessitate the use of existing mature build systems.