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


def cmd_get(args):
    device = open_device()
    key_cw, key_ccw, divider = get_config(device)
    print(f"key_cw   = 0x{key_cw:04X}")
    print(f"key_ccw  = 0x{key_ccw:04X}")
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
    p_set.add_argument("--cw", type=parse_int, help="clockwise key (hex, e.g. 0xE9)")
    p_set.add_argument(
        "--ccw", type=parse_int, help="counter-clockwise key (hex, e.g. 0xEA)"
    )
    p_set.add_argument("--divider", type=parse_int, help="encoder divider")

    sub.add_parser("save", help="persist current config to flash")
    sub.add_parser("load", help="reload config from flash")
    sub.add_parser("defaults", help="reset to compiled-in defaults")

    args = parser.parse_args()
    commands = {
        "get": cmd_get,
        "set": cmd_set,
        "save": cmd_save,
        "load": cmd_load,
        "defaults": cmd_defaults,
    }
    commands[args.command](args)


if __name__ == "__main__":
    main()
