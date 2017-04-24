all: main

main: assets
	go build

test: main glide
	go test `glide novendor`

assets: bindata
	go-bindata -o docstring/assets.go -pkg="docstring" \
	-prefix="docstring/" docstring/*.md

bindata:
	go get github.com/jteeuwen/go-bindata/...

glide:
	go get github.com/Masterminds/glide

init: glide bindata
	glide install
