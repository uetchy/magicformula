default: build
build:
	go-bindata templates/
	go fmt
	go build
run: build
	./magicformula $(ARGS)
