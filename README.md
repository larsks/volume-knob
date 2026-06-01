# Volume knob

Use a Raspberry Pi Pico to turn a quadrature encoder into a volume knob for your computer. Works great with something like <https://a.co/d/04KHyvDu>:

![Picture of a quadrature encoder](encoder.jpg)

## Hardware

You will need a [Raspberry Pi Pico] or something similar (I use a [Waveshare RP2040-Zero]).

[raspberry pi pico]: https://www.waveshare.com/wiki/RP2040-Zero
[waveshare RP2040-Zero]: https://www.waveshare.com/wiki/RP2040-Zero

## Configuration

In [config.h](config.h), you may want to set:

- `ENCODER_PIN_A`, `ENCODER_PIN_B`

  Defines the GPIOs to which the encoder outputs are attached.

- `KEY_CW`, `KEY_CCW`

  Defines which key events are sent on clockwise and counter-clockwise rotation. Defaults to volume up/volume down.

- `ENCODER_DIVIDER`

  Defines how many events we need to receive from the encoder before generating a key event.

## Building

This is a [pico-sdk]-based project. You will need a copy of the [pico-sdk], and you will need to set the `PICO_SDK_PATH` environment variable to point to the location of the sdk directory.

[pico-sdk]: https://github.com/raspberrypi/pico-sdk

## Updating the device

After building new code:

1. `picotool info -f --vid 0x1209 --pid 0x2641`

    This will reboot the Pico into bootsel mode. The command will fail like this:

    ```
    Tracking device serial number E66138935F4C6724 for reboot
    The device was asked to reboot into BOOTSEL mode so the command can be executed.
    Waiting for device to reboot.........

    Despite the reboot attempt, no accessible RP-series devices in BOOTSEL mode
    were found found with serial number E66138935F4C6724. It is possible the
    device is not responding, and will have to be manually entered into BOOTSEL
    mode.
    ```

    That's because the VID/PID change when it reboots. You can verify that it was successful by using your kernel log, which should show something like:

    ```
    kernel: usb 1-8: new full-speed USB device number 47 using xhci_hcd
    kernel: usb 1-8: New USB device found, idVendor=2e8a, idProduct=0003, bcdDevice= 1.00
    kernel: usb 1-8: New USB device strings: Mfr=1, Product=2, SerialNumber=3
    kernel: usb 1-8: Product: RP2 Boot
    kernel: usb 1-8: Manufacturer: Raspberry Pi
    kernel: usb 1-8: SerialNumber: E0C9125B0D9B
    ```

1. `picotool load build/volume_knob.uf2 -x`

    This will write the code to the device and reboot it.

## License

volume-knob -- Use a quadrature encoder as a volume knob\
Copyright (C) 2026 Lars Kellogg-Stedman

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
