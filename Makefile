.PHONY: build
build:
	go build -o ./out/autotiler .

.PHONY: run
run:
	go run . ./examples/2x3_packed.png output.local.png