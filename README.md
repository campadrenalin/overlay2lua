overlay2lua
===========

Covert PNG overlay files into positional matrices.

## Usage

    $ overlay2lua path/to/config.json path/to/character.png

This will output path/to/character.lua

## Build instructions

Install Go and its compilers on your machine first of all. If you've ever built
tmx2lua, you've already done this step. Then go through the following commands
to download the source of overlay2lua, install its dependencies, and build.

    $ git clone git@github.com:campadrenalin/overlay2lua.git
    $ cd overlay2lua
    $ go get github.com/bitly/go-simplejson # JSON lib dependency for config
    $ PLATFORM=linux # See below for available platforms
    $ make build/$PLATFORM/overlay2lua

You can currently compile for the following platforms:

 * linux
 * linux64
 * osx
 * windows32
 * windows64

To install, build the executable using the above instructions, and then
copy it to a location in your path, like ~/bin. No automatic installation
target is currently supported in the makefile.

### Note for Windows users

You're going to need to call the "make" line a bit differently, adding ".exe"
to the end of it. For example:

     $ make build/windows64/overlay2lua.exe

Versus, for example, the OSX build command:

    $ make build/osx/overlay2lua

## Precompiled binaries

The author of this software is not currently equipped to build for targets
other than 32-bit Linux. Other precompiled builds will be available if we
start using britta-bot to do automated builds and uploads - until then,
it's all manual. Sorry folks, I'm well aware what a pain in the butt that is.
