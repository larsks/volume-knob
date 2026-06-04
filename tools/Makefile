BINARY = vkcfg
BUILDFLAGS =
GEN_CONFIG ?= ../build/gen_config

$(BINARY): config_generated.go $(wildcard *.go) go.mod go.sum
	go build $(BUILDFLAGS) -o $@ .

config_generated.go: $(GEN_CONFIG) cmd/gen-config/main.go
	$(GEN_CONFIG) | go run ./cmd/gen-config > $@

clean:
	rm -f $(BINARY) config_generated.go

.PHONY: clean
