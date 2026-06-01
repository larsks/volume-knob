cmake_minimum_required(VERSION 3.13)

set(PICO_SDK_PATH "$ENV{HOME}/src/pico-sdk")
set(PICO_BOARD waveshare_rp2040_zero)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include(${PICO_SDK_PATH}/external/pico_sdk_import.cmake)

project(volume_knob C CXX ASM)
set(CMAKE_C_STANDARD 11)

pico_sdk_init()

add_executable(volume_knob
    main.c
    persist.c
    usb_descriptors.c
    version.c
)

target_include_directories(volume_knob PRIVATE ${CMAKE_CURRENT_SOURCE_DIR})

target_compile_definitions(volume_knob PRIVATE
    PICO_STDIO_USB_ENABLE_RESET_VIA_VENDOR_INTERFACE=1
    PICO_STDIO_USB_ENABLE_RESET_VIA_BAUD_RATE=1
    PICO_STDIO_USB_RESET_INTERFACE_SUPPORT_RESET_TO_BOOTSEL=1
    PICO_STDIO_USB_RESET_INTERFACE_SUPPORT_RESET_TO_FLASH_BOOT=1
    PICO_STDIO_USB_RESET_BOOTSEL_INTERFACE_DISABLE_MASK=0
    PICO_STDIO_USB_RESET_RESET_TO_FLASH_DELAY_MS=100
    PICO_STDIO_USB_RESET_MAGIC_BAUD_RATE=1200
)

pico_enable_stdio_usb(volume_knob 1)

target_link_libraries(volume_knob
    pico_stdlib
    hardware_gpio
    hardware_flash
    hardware_sync
    tinyusb_device
    tinyusb_board
)

target_link_options(volume_knob PRIVATE "LINKER:--undefined=version_info")

pico_add_extra_outputs(volume_knob)
