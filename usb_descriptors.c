#include "pico/unique_id.h"
#include "tusb.h" // IWYU pragma: keep

#include "usb_hid.h"

#define USB_VID 0x1209 // USB vendor id (https://pid.codes/1209/)
#define USB_PID 0x2641 // USB product id (https://pid.codes/1209/2641/)
#define USB_BCD 0x0200

tusb_desc_device_t const desc_device = {
    .bLength = sizeof(tusb_desc_device_t),
    .bDescriptorType = TUSB_DESC_DEVICE,
    .bcdUSB = USB_BCD,
    .bDeviceClass = 0,
    .bDeviceSubClass = 0,
    .bDeviceProtocol = 0,
    .bMaxPacketSize0 = CFG_TUD_ENDPOINT0_SIZE,
    .idVendor = USB_VID,
    .idProduct = USB_PID,
    .bcdDevice = 0x0100,
    .iManufacturer = 0x01,
    .iProduct = 0x02,
    .iSerialNumber = 0x03,
    .bNumConfigurations = 0x01,
};

uint8_t const *tud_descriptor_device_cb(void) {
  return (uint8_t const *)&desc_device;
}

uint8_t const desc_hid_report_consumer[] = {
    TUD_HID_REPORT_DESC_CONSUMER(HID_REPORT_ID(REPORT_ID_CONSUMER_CONTROL)),
    TUD_HID_REPORT_DESC_KEYBOARD(HID_REPORT_ID(REPORT_ID_KEYBOARD)),
};

// clang-format off
uint8_t const desc_hid_report_config[] = {
    0x06, 0x00, 0xFF,              // Usage Page (Vendor Defined 0xFF00)
    0x09, 0x01,                    // Usage (Vendor Usage 1)
    0xA1, 0x01,                    // Collection (Application)

    // Report ID 1: Config Data (feature report, 8 bytes)
    0x85, REPORT_ID_CONFIG,        //   Report ID
    0x09, 0x02,                    //   Usage (Vendor Usage 2)
    0x15, 0x00,                    //   Logical Minimum (0)
    0x26, 0xFF, 0x00,              //   Logical Maximum (255)
    0x75, 0x08,                    //   Report Size (8)
    0x95, 0x0A,                    //   Report Count (10)
    0xB1, 0x02,                    //   Feature (Data, Variable, Absolute)

    // Report ID 2: Command (feature report, 1x uint8)
    0x85, REPORT_ID_COMMAND,       //   Report ID
    0x09, 0x03,                    //   Usage (Vendor Usage 3)
    0x15, 0x00,                    //   Logical Minimum (0)
    0x25, 0x04,                    //   Logical Maximum (4)
    0x75, 0x08,                    //   Report Size (8)
    0x95, 0x01,                    //   Report Count (1)
    0xB1, 0x02,                    //   Feature (Data, Variable, Absolute)

    0xC0,                          // End Collection
};
// clang-format on

uint8_t const *tud_hid_descriptor_report_cb(uint8_t instance) {
  if (instance == HID_INSTANCE_CONFIG)
    return desc_hid_report_config;
  return desc_hid_report_consumer;
}

enum {
  ITF_NUM_HID,
  ITF_NUM_HID_CONFIG,
  ITF_NUM_TOTAL,
};

#define CONFIG_TOTAL_LEN (TUD_CONFIG_DESC_LEN + 2 * TUD_HID_DESC_LEN)

#define EPNUM_HID 0x81
#define EPNUM_HID_CONFIG 0x82

uint8_t const desc_configuration[] = {
    TUD_CONFIG_DESCRIPTOR(1, ITF_NUM_TOTAL, 0, CONFIG_TOTAL_LEN,
                          TUSB_DESC_CONFIG_ATT_REMOTE_WAKEUP, 100),

    TUD_HID_DESCRIPTOR(ITF_NUM_HID, 0, HID_ITF_PROTOCOL_NONE,
                       sizeof(desc_hid_report_consumer), EPNUM_HID,
                       CFG_TUD_HID_EP_BUFSIZE, 5),

    TUD_HID_DESCRIPTOR(ITF_NUM_HID_CONFIG, 4, HID_ITF_PROTOCOL_NONE,
                       sizeof(desc_hid_report_config), EPNUM_HID_CONFIG,
                       CFG_TUD_HID_EP_BUFSIZE, 5),
};

uint8_t const *tud_descriptor_configuration_cb(uint8_t index) {
  (void)index;
  return desc_configuration;
}

static char serial_str[PICO_UNIQUE_BOARD_ID_SIZE_BYTES * 2 + 1];

char const *string_desc_arr[] = {
    [0] = (const char[]){0x09, 0x04},
    [1] = "The Odd Bit",
    [2] = "Volume Knob",
    [3] = serial_str,
    [4] = "Config",
};

static uint16_t _desc_str[32 + 1];

uint16_t const *tud_descriptor_string_cb(uint8_t index, uint16_t langid) {
  (void)langid;
  size_t chr_count;

  if (!serial_str[0]) {
    pico_get_unique_board_id_string(serial_str, sizeof(serial_str));
  }

  switch (index) {
  case 0:
    memcpy(&_desc_str[1], string_desc_arr[0], 2);
    chr_count = 1;
    break;
  default:
    if (index >= sizeof(string_desc_arr) / sizeof(string_desc_arr[0]))
      return NULL;

    const char *str = string_desc_arr[index];
    chr_count = strlen(str);
    size_t const max_count = sizeof(_desc_str) / sizeof(_desc_str[0]) - 1;
    if (chr_count > max_count)
      chr_count = max_count;

    for (size_t i = 0; i < chr_count; i++) {
      _desc_str[1 + i] = str[i];
    }
    break;
  }

  _desc_str[0] = (uint16_t)((TUSB_DESC_STRING << 8) | (2 * chr_count + 2));
  return _desc_str;
}
