.PHONY: all clean init run build_macos build_linux release_macos release_linux

GO_VERSION=1.10

MACOS_X64_BIN_TARGET=${CURDIR}/_build/release/macos_x64/serve
MACOS_X64_ZIP_TARGET=${CURDIR}/_build/release/macos_x64/serve-macos_x64.zip

LINUX_X64_BIN_TARGET=${CURDIR}/_build/release/linux_x64/serve
LINUX_X64_ZIP_TARGET=${CURDIR}/_build/release/linux_x64/serve-linux_x64.zip

all:
	clean

clean:
	rm -rf _build
	rm -f ./serve

init:
	go mod vendor

run:
	go run serve.go

build_macos:
	go build -o ${MACOS_X64_BIN_TARGET} .
	ln -s ${MACOS_X64_BIN_TARGET} ${CURDIR}/serve

build_linux:
	docker pull golang:${GO_VERSION}
	docker run \
		--rm -v ${CURDIR}:/go/src/github.com/sandeepraju/serve \
		-w /go/src/github.com/sandeepraju/serve golang:${GO_VERSION} \
		sh -c 'go vendor && go build -o ./_build/release/linux_x64/serve .'

release_macos:
	zip -j ${MACOS_X64_ZIP_TARGET} ${MACOS_X64_BIN_TARGET}
	shasum -a256 ${MACOS_X64_ZIP_TARGET}
	find _build -name *.zip

release_linux:
	zip -j ${LINUX_X64_ZIP_TARGET} ${LINUX_X64_BIN_TARGET}
	shasum -a256 ${LINUX_X64_ZIP_TARGET}
	find _build -name *.zip
