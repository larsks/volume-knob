#!/usr/bin/env python3
"""Configuration tool for the Volume Knob."""

import argparse
import struct
import sys

import hid

VID = 0x1209
PID = 0x2641
VENDOR_USAGE_PAGE = 0xFF00
CONFIG_INTERFACE_NUMBER = 3

REPORT_ID_CONFIG = 1
REPORT_ID_COMMAND = 2

CONFIG_REPORT_SIZE = 6

CMD_SAVE = 1
CMD_LOAD = 2
CMD_DEFAULTS = 3

CONSUMER_KEYS = {
    "control": 0x0001,
    "power": 0x0030,
    "reset": 0x0031,
    "sleep": 0x0032,
    "brightness_increment": 0x006F,
    "brightness_decrement": 0x0070,
    "wireless_radio_controls": 0x000C,
    "wireless_radio_buttons": 0x00C6,
    "wireless_radio_led": 0x00C7,
    "wireless_radio_slider_switch": 0x00C8,
    "play_pause": 0x00CD,
    "scan_next": 0x00B5,
    "scan_previous": 0x00B6,
    "stop": 0x00B7,
    "volume": 0x00E0,
    "mute": 0x00E2,
    "bass": 0x00E3,
    "treble": 0x00E4,
    "bass_boost": 0x00E5,
    "volume_increment": 0x00E9,
    "volume_decrement": 0x00EA,
    "bass_increment": 0x0152,
    "bass_decrement": 0x0153,
    "treble_increment": 0x0154,
    "treble_decrement": 0x0155,
    "al_consumer_control_configuration": 0x0183,
    "al_email_reader": 0x018A,
    "al_calculator": 0x0192,
    "al_local_browser": 0x0194,
    "ac_search": 0x0221,
    "ac_home": 0x0223,
    "ac_back": 0x0224,
    "ac_forward": 0x0225,
    "ac_stop": 0x0226,
    "ac_refresh": 0x0227,
    "ac_bookmarks": 0x022A,
    "ac_pan": 0x0238,
}

CONSUMER_KEY_NAMES = {v: k for k, v in CONSUMER_KEYS.items()}


def key_name(code):
    name = CONSUMER_KEY_NAMES.get(code)
    if name:
        return f"0x{code:04X} ({name})"
    return f"0x{code:04X}"


def parse_key(s):
    name = s.lower()
    if name in CONSUMER_KEYS:
        return CONSUMER_KEYS[name]
    try:
        return int(s, 0)
    except ValueError:
        raise argparse.ArgumentTypeError(
            f"unknown key name {s!r} (use 'list-keys' to see valid names)"
        )


def find_config_interface():
    for dev in hid.enumerate(VID, PID):
        if dev["usage_page"] == VENDOR_USAGE_PAGE:
            return dev["path"]
    # libusb backend doesn't populate usage_page; fall back to interface number
    for dev in hid.enumerate(VID, PID):
        if dev["interface_number"] == CONFIG_INTERFACE_NUMBER:
            return dev["path"]
    return None


def open_device():
    path = find_config_interface()
    if path is None:
        print("Error: Volume Knob not found", file=sys.stderr)
        sys.exit(1)
    device = hid.device()
    device.open_path(path)
    return device


def get_config(device):
    data = device.get_feature_report(REPORT_ID_CONFIG, CONFIG_REPORT_SIZE + 1)
    key_cw, key_ccw, divider = struct.unpack_from("<HHH", bytes(data), 1)
    return key_cw, key_ccw, divider


def set_config(device, key_cw, key_ccw, divider):
    data = struct.pack("<BHHH", REPORT_ID_CONFIG, key_cw, key_ccw, divider)
    device.send_feature_report(data)


def send_command(device, cmd):
    data = struct.pack("<BB", REPORT_ID_COMMAND, cmd)
    device.send_feature_report(data)


def parse_int(s):
    return int(s, 0)


def cmd_list_keys(args):
    for name in sorted(CONSUMER_KEYS):
        print(f"  {name:40s} 0x{CONSUMER_KEYS[name]:04X}")


def cmd_get(args):
    device = open_device()
    key_cw, key_ccw, divider = get_config(device)
    print(f"key_cw   = {key_name(key_cw)}")
    print(f"key_ccw  = {key_name(key_ccw)}")
    print(f"divider  = {divider}")


def cmd_set(args):
    device = open_device()
    key_cw, key_ccw, divider = get_config(device)
    if args.cw is not None:
        key_cw = args.cw
    if args.ccw is not None:
        key_ccw = args.ccw
    if args.divider is not None:
        divider = args.divider
    set_config(device, key_cw, key_ccw, divider)
    print("OK")


def cmd_save(args):
    device = open_device()
    send_command(device, CMD_SAVE)
    print("OK")


def cmd_load(args):
    device = open_device()
    send_command(device, CMD_LOAD)
    print("OK")


def cmd_defaults(args):
    device = open_device()
    send_command(device, CMD_DEFAULTS)
    print("OK")


def main():
    parser = argparse.ArgumentParser(description="Volume Knob configuration")
    sub = parser.add_subparsers(dest="command", required=True)

    sub.add_parser("get", help="show current configuration")

    p_set = sub.add_parser("set", help="update configuration values")
    p_set.add_argument(
        "--cw", type=parse_key, help="clockwise key (name or hex, e.g. volume_increment or 0xE9)"
    )
    p_set.add_argument(
        "--ccw", type=parse_key, help="counter-clockwise key (name or hex, e.g. volume_decrement or 0xEA)"
    )
    p_set.add_argument("--divider", type=parse_int, help="encoder divider")

    sub.add_parser("save", help="persist current config to flash")
    sub.add_parser("load", help="reload config from flash")
    sub.add_parser("defaults", help="reset to compiled-in defaults")
    sub.add_parser("list-keys", help="list known key names")

    args = parser.parse_args()
    commands = {
        "get": cmd_get,
        "set": cmd_set,
        "save": cmd_save,
        "load": cmd_load,
        "defaults": cmd_defaults,
        "list-keys": cmd_list_keys,
    }
    commands[args.command](args)


if __name__ == "__main__":
    main()
