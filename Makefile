all: install

install: build test
	go install github.com/kode4food/sputter/cmd/sputter

test: build dep-lint
	golint ./...
	go vet ./...
	go test ./...

build: dependencies assets

assets: dep-snapshot
	go-snapshot -pkg assets -out assets/assets.go docstring/*.md core/*.lisp

dependencies: dep-dep
	dep ensure

dep-dep:
	go get github.com/golang/dep/cmd/dep

dep-snapshot:
	go get github.com/kode4food/go-snapshot

dep-lint:
	go get github.com/golang/lint/golint

upgrade-deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/kode4food/go-snapshot
	go get -u github.com/golang/lint/golint
