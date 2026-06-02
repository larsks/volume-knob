BINARY = vkcfg
BUILDFLAGS =

$(BINARY): $(wildcard *.go) go.mod go.sum
	go build $(BUILDFLAGS) -o $@ .

clean:
	rm -f $(BINARY)

.PHONY: clean
