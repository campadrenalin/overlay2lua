.PHONY = clean fmt tar

all: fmt build/osx/overlay2lua build/windows32/overlay2lua.exe build/windows64/overlay2lua.exe build/linux/overlay2lua build/linux64/overlay2lua

tar: all build/overlay2lua.windows32.tar.gz build/overlay2lua.windows64.tar.gz build/overlay2lua.osx.tar.gz build/overlay2lua.linux.tar.gz build/overlay2lua.linux64.tar.gz

build/osx/overlay2lua: overlay2lua.go
	mkdir -p build/osx
	GOOS=darwin GOARCH=amd64 go build -o $@

build/linux64/overlay2lua: overlay2lua.go
	mkdir -p build/linux64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@

build/linux/overlay2lua: overlay2lua.go
	mkdir -p build/linux
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $@

build/windows64/overlay2lua.exe: overlay2lua.go
	mkdir -p build/windows64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $@

build/windows32/overlay2lua.exe: overlay2lua.go
	mkdir -p build/windows32
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o $@

build/overlay2lua.windows%.tar.gz: build/windows%/overlay2lua.exe
	cd build/windows$* && tar -cf overlay2lua.windows$*.tar.gz overlay2lua.exe
	mv build/windows$*/overlay2lua.windows$*.tar.gz build

build/overlay2lua.%.tar.gz: build/%/overlay2lua
	cd build/$* && tar -cf overlay2lua.$*.tar.gz overlay2lua
	mv build/$*/overlay2lua.$*.tar.gz build

fmt: 
	go fmt overlay2lua.go

clean: 
	go clean
	rm -rf build
