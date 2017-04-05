all: main

main: assets
	go build

assets:
	go-bindata -o docstring/assets.go -pkg="docstring" \
	-prefix="docstring/" docstring/*.md
