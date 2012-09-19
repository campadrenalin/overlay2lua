overlay2lua
===========

Covert PNG overlay files into positional matrices.

## Usage

    $ overlay2lua path/to/config.json path/to/character.png

This will output path/to/character.lua

## Build instructions

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
