all: main

main: assets
	go build github.com/kode4food/sputter/cmd/sputter

test: main glide
	golint `glide novendor`
	go vet `glide novendor`
	go test `glide novendor`

assets: snapshot
	go-snapshot -pkg assets -out assets/assets.go docstring/*.md core/*.lisp

glide:
	go get github.com/Masterminds/glide

snapshot:
	go get github.com/kode4food/go-snapshot

lint:
	go get github.com/golang/lint/golint

init: glide assets lint
	glide install

install: init main test
	go install github.com/kode4food/sputter/cmd/sputter

upgrade-deps:
	go get -u github.com/Masterminds/glide
	go get -u github.com/kode4food/go-snapshot
	go get -u github.com/golang/lint/golint
