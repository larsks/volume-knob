BINARY = vkcfg
BUILDFLAGS = -tags libusb

$(BINARY): $(wildcard *.go) go.mod go.sum
	go build $(BUILDFLAGS) -o $@ .

clean:
	rm -f $(BINARY)

.PHONY: clean
