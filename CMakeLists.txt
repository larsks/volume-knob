cmake_minimum_required(VERSION 3.13)

set(PICO_SDK_PATH "$ENV{HOME}/src/pico-sdk")
set(PICO_BOARD waveshare_rp2040_zero)

include(${PICO_SDK_PATH}/external/pico_sdk_import.cmake)

project(volume_knob C CXX ASM)
set(CMAKE_C_STANDARD 11)

pico_sdk_init()

add_executable(volume_knob
    main.c
    usb_descriptors.c
)

target_include_directories(volume_knob PRIVATE ${CMAKE_CURRENT_SOURCE_DIR})

target_link_libraries(volume_knob
    pico_stdlib
    hardware_gpio
    tinyusb_device
    tinyusb_board
)

pico_add_extra_outputs(volume_knob)
