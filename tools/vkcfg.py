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

CONFIG_REPORT_SIZE = 8

CMD_SAVE = 1
CMD_LOAD = 2
CMD_DEFAULTS = 3

KEY_TYPE_CONSUMER = 0
KEY_TYPE_KEYBOARD = 1

# fmt: off
ALL_KEYS = {
    # Consumer control keys (from TinyUSB HID_USAGE_CONSUMER_*)
    "consumer_control": (KEY_TYPE_CONSUMER, 0x0001),
    "power": (KEY_TYPE_CONSUMER, 0x0030),
    "reset": (KEY_TYPE_CONSUMER, 0x0031),
    "sleep": (KEY_TYPE_CONSUMER, 0x0032),
    "brightness_increment": (KEY_TYPE_CONSUMER, 0x006F),
    "brightness_decrement": (KEY_TYPE_CONSUMER, 0x0070),
    "wireless_radio_controls": (KEY_TYPE_CONSUMER, 0x000C),
    "wireless_radio_buttons": (KEY_TYPE_CONSUMER, 0x00C6),
    "wireless_radio_led": (KEY_TYPE_CONSUMER, 0x00C7),
    "wireless_radio_slider_switch": (KEY_TYPE_CONSUMER, 0x00C8),
    "play_pause": (KEY_TYPE_CONSUMER, 0x00CD),
    "scan_next": (KEY_TYPE_CONSUMER, 0x00B5),
    "scan_previous": (KEY_TYPE_CONSUMER, 0x00B6),
    "stop": (KEY_TYPE_CONSUMER, 0x00B7),
    "volume": (KEY_TYPE_CONSUMER, 0x00E0),
    "mute": (KEY_TYPE_CONSUMER, 0x00E2),
    "bass": (KEY_TYPE_CONSUMER, 0x00E3),
    "treble": (KEY_TYPE_CONSUMER, 0x00E4),
    "bass_boost": (KEY_TYPE_CONSUMER, 0x00E5),
    "volume_increment": (KEY_TYPE_CONSUMER, 0x00E9),
    "volume_decrement": (KEY_TYPE_CONSUMER, 0x00EA),
    "bass_increment": (KEY_TYPE_CONSUMER, 0x0152),
    "bass_decrement": (KEY_TYPE_CONSUMER, 0x0153),
    "treble_increment": (KEY_TYPE_CONSUMER, 0x0154),
    "treble_decrement": (KEY_TYPE_CONSUMER, 0x0155),
    "al_consumer_control_configuration": (KEY_TYPE_CONSUMER, 0x0183),
    "al_email_reader": (KEY_TYPE_CONSUMER, 0x018A),
    "al_calculator": (KEY_TYPE_CONSUMER, 0x0192),
    "al_local_browser": (KEY_TYPE_CONSUMER, 0x0194),
    "ac_search": (KEY_TYPE_CONSUMER, 0x0221),
    "ac_home": (KEY_TYPE_CONSUMER, 0x0223),
    "ac_back": (KEY_TYPE_CONSUMER, 0x0224),
    "ac_forward": (KEY_TYPE_CONSUMER, 0x0225),
    "ac_stop": (KEY_TYPE_CONSUMER, 0x0226),
    "ac_refresh": (KEY_TYPE_CONSUMER, 0x0227),
    "ac_bookmarks": (KEY_TYPE_CONSUMER, 0x022A),
    "ac_pan": (KEY_TYPE_CONSUMER, 0x0238),

    # Keyboard keys (from TinyUSB HID_KEY_*)
    "a": (KEY_TYPE_KEYBOARD, 0x04),
    "b": (KEY_TYPE_KEYBOARD, 0x05),
    "c": (KEY_TYPE_KEYBOARD, 0x06),
    "d": (KEY_TYPE_KEYBOARD, 0x07),
    "e": (KEY_TYPE_KEYBOARD, 0x08),
    "f": (KEY_TYPE_KEYBOARD, 0x09),
    "g": (KEY_TYPE_KEYBOARD, 0x0A),
    "h": (KEY_TYPE_KEYBOARD, 0x0B),
    "i": (KEY_TYPE_KEYBOARD, 0x0C),
    "j": (KEY_TYPE_KEYBOARD, 0x0D),
    "k": (KEY_TYPE_KEYBOARD, 0x0E),
    "l": (KEY_TYPE_KEYBOARD, 0x0F),
    "m": (KEY_TYPE_KEYBOARD, 0x10),
    "n": (KEY_TYPE_KEYBOARD, 0x11),
    "o": (KEY_TYPE_KEYBOARD, 0x12),
    "p": (KEY_TYPE_KEYBOARD, 0x13),
    "q": (KEY_TYPE_KEYBOARD, 0x14),
    "r": (KEY_TYPE_KEYBOARD, 0x15),
    "s": (KEY_TYPE_KEYBOARD, 0x16),
    "t": (KEY_TYPE_KEYBOARD, 0x17),
    "u": (KEY_TYPE_KEYBOARD, 0x18),
    "v": (KEY_TYPE_KEYBOARD, 0x19),
    "w": (KEY_TYPE_KEYBOARD, 0x1A),
    "x": (KEY_TYPE_KEYBOARD, 0x1B),
    "y": (KEY_TYPE_KEYBOARD, 0x1C),
    "z": (KEY_TYPE_KEYBOARD, 0x1D),
    "1": (KEY_TYPE_KEYBOARD, 0x1E),
    "2": (KEY_TYPE_KEYBOARD, 0x1F),
    "3": (KEY_TYPE_KEYBOARD, 0x20),
    "4": (KEY_TYPE_KEYBOARD, 0x21),
    "5": (KEY_TYPE_KEYBOARD, 0x22),
    "6": (KEY_TYPE_KEYBOARD, 0x23),
    "7": (KEY_TYPE_KEYBOARD, 0x24),
    "8": (KEY_TYPE_KEYBOARD, 0x25),
    "9": (KEY_TYPE_KEYBOARD, 0x26),
    "0": (KEY_TYPE_KEYBOARD, 0x27),
    "enter": (KEY_TYPE_KEYBOARD, 0x28),
    "escape": (KEY_TYPE_KEYBOARD, 0x29),
    "backspace": (KEY_TYPE_KEYBOARD, 0x2A),
    "tab": (KEY_TYPE_KEYBOARD, 0x2B),
    "space": (KEY_TYPE_KEYBOARD, 0x2C),
    "minus": (KEY_TYPE_KEYBOARD, 0x2D),
    "equal": (KEY_TYPE_KEYBOARD, 0x2E),
    "bracket_left": (KEY_TYPE_KEYBOARD, 0x2F),
    "bracket_right": (KEY_TYPE_KEYBOARD, 0x30),
    "backslash": (KEY_TYPE_KEYBOARD, 0x31),
    "semicolon": (KEY_TYPE_KEYBOARD, 0x33),
    "apostrophe": (KEY_TYPE_KEYBOARD, 0x34),
    "grave": (KEY_TYPE_KEYBOARD, 0x35),
    "comma": (KEY_TYPE_KEYBOARD, 0x36),
    "period": (KEY_TYPE_KEYBOARD, 0x37),
    "slash": (KEY_TYPE_KEYBOARD, 0x38),
    "caps_lock": (KEY_TYPE_KEYBOARD, 0x39),
    "f1": (KEY_TYPE_KEYBOARD, 0x3A),
    "f2": (KEY_TYPE_KEYBOARD, 0x3B),
    "f3": (KEY_TYPE_KEYBOARD, 0x3C),
    "f4": (KEY_TYPE_KEYBOARD, 0x3D),
    "f5": (KEY_TYPE_KEYBOARD, 0x3E),
    "f6": (KEY_TYPE_KEYBOARD, 0x3F),
    "f7": (KEY_TYPE_KEYBOARD, 0x40),
    "f8": (KEY_TYPE_KEYBOARD, 0x41),
    "f9": (KEY_TYPE_KEYBOARD, 0x42),
    "f10": (KEY_TYPE_KEYBOARD, 0x43),
    "f11": (KEY_TYPE_KEYBOARD, 0x44),
    "f12": (KEY_TYPE_KEYBOARD, 0x45),
    "print_screen": (KEY_TYPE_KEYBOARD, 0x46),
    "scroll_lock": (KEY_TYPE_KEYBOARD, 0x47),
    "pause": (KEY_TYPE_KEYBOARD, 0x48),
    "insert": (KEY_TYPE_KEYBOARD, 0x49),
    "home": (KEY_TYPE_KEYBOARD, 0x4A),
    "page_up": (KEY_TYPE_KEYBOARD, 0x4B),
    "delete": (KEY_TYPE_KEYBOARD, 0x4C),
    "end": (KEY_TYPE_KEYBOARD, 0x4D),
    "page_down": (KEY_TYPE_KEYBOARD, 0x4E),
    "arrow_right": (KEY_TYPE_KEYBOARD, 0x4F),
    "arrow_left": (KEY_TYPE_KEYBOARD, 0x50),
    "arrow_down": (KEY_TYPE_KEYBOARD, 0x51),
    "arrow_up": (KEY_TYPE_KEYBOARD, 0x52),
    "num_lock": (KEY_TYPE_KEYBOARD, 0x53),
    "keypad_divide": (KEY_TYPE_KEYBOARD, 0x54),
    "keypad_multiply": (KEY_TYPE_KEYBOARD, 0x55),
    "keypad_subtract": (KEY_TYPE_KEYBOARD, 0x56),
    "keypad_add": (KEY_TYPE_KEYBOARD, 0x57),
    "keypad_enter": (KEY_TYPE_KEYBOARD, 0x58),
    "keypad_1": (KEY_TYPE_KEYBOARD, 0x59),
    "keypad_2": (KEY_TYPE_KEYBOARD, 0x5A),
    "keypad_3": (KEY_TYPE_KEYBOARD, 0x5B),
    "keypad_4": (KEY_TYPE_KEYBOARD, 0x5C),
    "keypad_5": (KEY_TYPE_KEYBOARD, 0x5D),
    "keypad_6": (KEY_TYPE_KEYBOARD, 0x5E),
    "keypad_7": (KEY_TYPE_KEYBOARD, 0x5F),
    "keypad_8": (KEY_TYPE_KEYBOARD, 0x60),
    "keypad_9": (KEY_TYPE_KEYBOARD, 0x61),
    "keypad_0": (KEY_TYPE_KEYBOARD, 0x62),
    "keypad_decimal": (KEY_TYPE_KEYBOARD, 0x63),
    "application": (KEY_TYPE_KEYBOARD, 0x65),
    "key_power": (KEY_TYPE_KEYBOARD, 0x66),
    "keypad_equal": (KEY_TYPE_KEYBOARD, 0x67),
    "f13": (KEY_TYPE_KEYBOARD, 0x68),
    "f14": (KEY_TYPE_KEYBOARD, 0x69),
    "f15": (KEY_TYPE_KEYBOARD, 0x6A),
    "f16": (KEY_TYPE_KEYBOARD, 0x6B),
    "f17": (KEY_TYPE_KEYBOARD, 0x6C),
    "f18": (KEY_TYPE_KEYBOARD, 0x6D),
    "f19": (KEY_TYPE_KEYBOARD, 0x6E),
    "f20": (KEY_TYPE_KEYBOARD, 0x6F),
    "f21": (KEY_TYPE_KEYBOARD, 0x70),
    "f22": (KEY_TYPE_KEYBOARD, 0x71),
    "f23": (KEY_TYPE_KEYBOARD, 0x72),
    "f24": (KEY_TYPE_KEYBOARD, 0x73),
    "execute": (KEY_TYPE_KEYBOARD, 0x74),
    "help": (KEY_TYPE_KEYBOARD, 0x75),
    "menu": (KEY_TYPE_KEYBOARD, 0x76),
    "select": (KEY_TYPE_KEYBOARD, 0x77),
    "key_stop": (KEY_TYPE_KEYBOARD, 0x78),
    "again": (KEY_TYPE_KEYBOARD, 0x79),
    "undo": (KEY_TYPE_KEYBOARD, 0x7A),
    "cut": (KEY_TYPE_KEYBOARD, 0x7B),
    "copy": (KEY_TYPE_KEYBOARD, 0x7C),
    "paste": (KEY_TYPE_KEYBOARD, 0x7D),
    "find": (KEY_TYPE_KEYBOARD, 0x7E),
    "key_mute": (KEY_TYPE_KEYBOARD, 0x7F),
    "volume_up": (KEY_TYPE_KEYBOARD, 0x80),
    "volume_down": (KEY_TYPE_KEYBOARD, 0x81),
    "control_left": (KEY_TYPE_KEYBOARD, 0xE0),
    "shift_left": (KEY_TYPE_KEYBOARD, 0xE1),
    "alt_left": (KEY_TYPE_KEYBOARD, 0xE2),
    "gui_left": (KEY_TYPE_KEYBOARD, 0xE3),
    "control_right": (KEY_TYPE_KEYBOARD, 0xE4),
    "shift_right": (KEY_TYPE_KEYBOARD, 0xE5),
    "alt_right": (KEY_TYPE_KEYBOARD, 0xE6),
    "gui_right": (KEY_TYPE_KEYBOARD, 0xE7),
}
# fmt: on

KEY_NAMES = {v: k for k, v in ALL_KEYS.items()}


def key_name(key_type, code):
    name = KEY_NAMES.get((key_type, code))
    if name:
        return name
    prefix = "consumer" if key_type == KEY_TYPE_CONSUMER else "keyboard"
    return f"{prefix}:0x{code:04X}"


def parse_key(s):
    name = s.lower()
    if name in ALL_KEYS:
        return ALL_KEYS[name]
    try:
        return (KEY_TYPE_CONSUMER, int(s, 0))
    except ValueError:
        raise argparse.ArgumentTypeError(
            f"unknown key name {s!r} (use 'list-keys' to see valid names)"
        )


def find_config_interface():
    for dev in hid.enumerate(VID, PID):
        if dev["usage_page"] == VENDOR_USAGE_PAGE:
            return dev["path"]
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
    type_cw, type_ccw, key_cw, key_ccw, divider = struct.unpack_from(
        "<BBHHH", bytes(data), 1
    )
    return type_cw, type_ccw, key_cw, key_ccw, divider


def set_config(device, type_cw, type_ccw, key_cw, key_ccw, divider):
    data = struct.pack(
        "<BBBHHH", REPORT_ID_CONFIG, type_cw, type_ccw, key_cw, key_ccw, divider
    )
    device.send_feature_report(data)


def send_command(device, cmd):
    data = struct.pack("<BB", REPORT_ID_COMMAND, cmd)
    device.send_feature_report(data)


def parse_int(s):
    return int(s, 0)


def cmd_list_keys(args):
    for name in sorted(ALL_KEYS):
        key_type, code = ALL_KEYS[name]
        kind = "consumer" if key_type == KEY_TYPE_CONSUMER else "keyboard"
        print(f"  {name:40s} {kind:10s} 0x{code:04X}")


def cmd_get(args):
    device = open_device()
    type_cw, type_ccw, kc_cw, kc_ccw, divider = get_config(device)
    print(f"key_cw   = {key_name(type_cw, kc_cw)}")
    print(f"key_ccw  = {key_name(type_ccw, kc_ccw)}")
    print(f"divider  = {divider}")


def cmd_set(args):
    device = open_device()
    type_cw, type_ccw, key_cw, key_ccw, divider = get_config(device)
    if args.cw is not None:
        type_cw, key_cw = args.cw
    if args.ccw is not None:
        type_ccw, key_ccw = args.ccw
    if args.divider is not None:
        divider = args.divider
    set_config(device, type_cw, type_ccw, key_cw, key_ccw, divider)
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
        "--cw",
        type=parse_key,
        help="clockwise key (e.g. volume_increment, a, f1, 0xE9)",
    )
    p_set.add_argument(
        "--ccw",
        type=parse_key,
        help="counter-clockwise key (e.g. volume_decrement, b, f2, 0xEA)",
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
