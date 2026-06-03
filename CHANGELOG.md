# Changelog

## v1.0.0 (2026-06-02)

### Feat

- Support modifier keys when generating keyboard events
- Add "bootsel" subcommand to vkcfg

## v0.3.0 (2026-06-02)

### Feat

- Use hidraw interface instead of libusb to talk to knob

### Fix

- Avoid more linter warnings
- Code cleanup and comments
- Remove erroneous capitalization

## v0.2.0 (2026-06-01)

### Feat

- Replace python vkcfg.py with Go compiled binary
- Extend runtime config to include arbitrary keyboard events
- Allow symbolic key names in vkcfg.py
- Allow runtime configuration of key and divider
- Enable software reboot into bootsel mode

### Fix

- Embed version information
- Update manufacturer string
- Collapse ENCODER_STEPS_PER_DETENT into ENCODER_DIVIDER
- Move REPORT_ID_CONSUMER_CONTROL out of config.h
- Avoid dropping events when busy
- Fix multiple definitions of REPORT_ID_CONSUMER_CONTROL
- Use proper vendor/product code
- Fix ENCODER_DIVIDER handling

## v0.1.0 (2026-05-31)
