BUILD_BIN=bin/pslpreview

build:
	go build -ldflags "-X main.Version=$$(git rev-parse --short HEAD)" -o $(BUILD_BIN)

clean:
	rm -f bin/$(BUILD_BIN)
