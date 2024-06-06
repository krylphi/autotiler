.PHONY: build
build:
	go build -o ./out/autotiler .

.PHONY: unpack-example
unpack-example:
	go run . ./examples/2x3_packed.png output.local.png