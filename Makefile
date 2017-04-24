all: main

main: assets
	go build

assets: bindata
	go-bindata -o docstring/assets.go -pkg="docstring" \
	-prefix="docstring/" docstring/*.md

bindata:
	go get github.com/jteeuwen/go-bindata/...
