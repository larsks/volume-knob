all: build/volume_knob.uf2 tools/vkcfg

build:
	cmake -B build

build/volume_knob.uf2 build/gen_config: build
	$(MAKE) -C build

tools/vkcfg: build/gen_config
	$(MAKE) -C tools

clean:
	$(MAKE) -C tools clean
	rm -rf build
