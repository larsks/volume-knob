BINARY = vkcfg
BUILDFLAGS =
GEN_CONFIG ?= ../build/gen_config

CONFIG_GENERATED = cmd/vkcfg/config_generated.go

$(BINARY): $(CONFIG_GENERATED) $(wildcard cmd/vkcfg/*.go) go.mod go.sum
	go build $(BUILDFLAGS) -o $@ ./cmd/vkcfg/

$(CONFIG_GENERATED): $(GEN_CONFIG) cmd/gen-config/main.go
	$(GEN_CONFIG) | go run ./cmd/gen-config > $@

clean:
	rm -f $(BINARY) $(CONFIG_GENERATED)

.PHONY: clean
