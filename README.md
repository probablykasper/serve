# Serve

A simple HTTP server for serving static files.

## Installtion

Using Homebrew:
```
brew install probablykasper/tap/serve
```

Installing from source:
```
go get github.com/probablykasper/serve
```

Download binaries
- [Linux x64](https://github.com/probablykasper/releases/latest/serve-linux_x64.zip)
- [macOS x64](https://github.com/probablykasper/releases/latest/serve-macos_x64.zip)

## Usage

```
Serve 1.0.0
  
USAGE:
    serve [dir] [options...]

EXAMPLE:
    serve -p 80 ./website

OPTIONS:
    --address value, -a value  The IP address or hostname of the interface (default: "0.0.0.0")
    --port value, -p value     The port to listen on (default: 2233)
    --verbose, -V              Log requests
    --help, -h                 Show help menu
    --version, -v              print the version
```

## Dev instructions

### Getting started
1. Install [Go](https://golang.org/)
2. Run `go mod vendor` to install dependencies
3. Start by running `make run` or `go run serve.go`

### Commands

Start
```
make run
```

Build & release the project (assuming you are developing on macOS)

```
make clean build_macos build_linux release_macos release_linux
```
