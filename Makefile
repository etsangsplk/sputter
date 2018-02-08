all: install

install: build test
	go install github.com/kode4food/sputter/cmd/sputter

test: build dep-lint
	golint `glide novendor`
	go vet `glide novendor`
	go test `glide novendor`

build: dependencies assets

assets: dep-snapshot
	go-snapshot -pkg assets -out assets/assets.go docstring/*.md core/*.lisp

dependencies: dep-glide
	glide install

dep-glide:
	go get github.com/Masterminds/glide

dep-snapshot:
	go get github.com/kode4food/go-snapshot

dep-lint:
	go get github.com/golang/lint/golint

upgrade-deps:
	go get -u github.com/Masterminds/glide
	go get -u github.com/kode4food/go-snapshot
	go get -u github.com/golang/lint/golint
