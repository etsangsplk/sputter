all:
	go-bindata -o docstring/docstring.go -pkg="docstring" -prefix="docstring/" docstring/*.md
