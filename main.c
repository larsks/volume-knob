#include "bsp/board_api.h"
#include "hardware/gpio.h"
#include "pico/bootrom.h"
#include "tusb.h"

#include <string.h>

#include "config.h"
#include "persist.h"
#include "usb_hid.h"

// Quadrature state transition table.
// Index: (prev_AB << 2) | curr_AB
// Values: 0 = invalid/no change, 1 = CW, -1 = CCW
// clang-format off
static const int8_t encoder_table[16] = {
  0,  -1, 1,   0,
  1,  0,  0,   -1,
  -1, 0,  0,  1,
  0, 1, -1, 0,
};
// clang-format on

static config_t config;
static uint8_t encoder_prev_state;
static int32_t encoder_accum;
static bool pending_release;

static uint8_t read_encoder_state(void) {
  uint8_t a = gpio_get(ENCODER_PIN_A) ? 1 : 0;
  uint8_t b = gpio_get(ENCODER_PIN_B) ? 1 : 0;
  return (a << 1) | b;
}

// Initialize the Pico GPIO pins that are connected to the encoder.
static void encoder_init(void) {
  gpio_init(ENCODER_PIN_A);
  gpio_set_dir(ENCODER_PIN_A, GPIO_IN);
  gpio_pull_up(ENCODER_PIN_A);

  gpio_init(ENCODER_PIN_B);
  gpio_set_dir(ENCODER_PIN_B, GPIO_IN);
  gpio_pull_up(ENCODER_PIN_B);

  encoder_prev_state = read_encoder_state();
  encoder_accum = 0;
}

// Send an HID key event. @type determins if we send a keyboard event
// or a consumer control event.
static bool send_key(uint8_t type, uint16_t code, uint8_t modifier) {
  if (!tud_hid_ready())
    return false;

  if (type == KEY_TYPE_KEYBOARD) {
    uint8_t report[8] = {0};
    report[0] = modifier;
    report[2] = (uint8_t)code;
    tud_hid_report(REPORT_ID_KEYBOARD, report, sizeof(report));
  } else {
    tud_hid_report(REPORT_ID_CONSUMER_CONTROL, &code, sizeof(code));
  }

  pending_release = true;
  return true;
}

// Read the current state of the encoder, determine in which direction it
// moved, and send the appropriate key down event. Key release is handled by
// tud_hid_report_complete_cb.
static void encoder_task(void) {
  uint8_t curr = read_encoder_state();

  // Create an index into the encoder table by combining
  // the previous encoder state and the current encoder state.
  // This is effectively creating a (x,y) tuple.
  uint8_t index = (encoder_prev_state << 2) | curr;
  encoder_prev_state = curr;

  int8_t delta = encoder_table[index];
  if (delta == 0)
    return;

  encoder_accum += delta;

  if (encoder_accum >= config.divider) {
    if (send_key(config.type_cw, config.key_cw, config.mod_cw))
      encoder_accum = 0;
  } else if (encoder_accum <= -config.divider) {
    if (send_key(config.type_ccw, config.key_ccw, config.mod_ccw))
      encoder_accum = 0;
  }
}

int main(void) {
  board_init();
  tusb_init();
  config_load(&config);
  encoder_init();

  while (1) {
    tud_task();
    encoder_task();
  }
}

// Triggered after an HID event has been successfully delivered to the host.
// In other words, this is called after a key down event has been sent
// successfully, so we generate the corresponding release event.
void tud_hid_report_complete_cb(uint8_t instance, uint8_t const *report,
                                uint16_t len) {

  // We only need a release event if the triggering report was a key event.
  // Both keyboard and consumer key events are sent through the
  // HID_INSTANCE_CONSUMER instance.
  if (instance != HID_INSTANCE_CONSUMER || !pending_release)
    return;

  pending_release = false;
  uint8_t release[8] = {0};
  tud_hid_n_report(instance, report[0], release, len - 1);
}

uint16_t tud_hid_get_report_cb(uint8_t instance, uint8_t report_id,
                               hid_report_type_t report_type, uint8_t *buffer,
                               uint16_t reqlen) {
  (void)report_type;
  (void)reqlen;

  // Read current configuration and return it via HID report.
  if (instance == HID_INSTANCE_CONFIG && report_id == REPORT_ID_CONFIG) {
    memcpy(buffer, &config, CONFIG_REPORT_SIZE);
    return CONFIG_REPORT_SIZE;
  }

  return 0;
}

void tud_hid_set_report_cb(uint8_t instance, uint8_t report_id,
                           hid_report_type_t report_type, uint8_t const *buffer,
                           uint16_t bufsize) {
  (void)report_type;

  if (instance != HID_INSTANCE_CONFIG)
    return;

  if (report_id == REPORT_ID_CONFIG && bufsize >= CONFIG_REPORT_SIZE) {
    // Receive new configuration via HID report and write it to runtime
    // configuration.
    memcpy(&config, buffer, CONFIG_REPORT_SIZE);
  } else if (report_id == REPORT_ID_COMMAND && bufsize >= COMMAND_REPORT_SIZE) {
    uint8_t cmd = buffer[0];
    if (cmd == CONFIG_CMD_SAVE)
      // Save runtime configuration to flash (persist across reboots).
      config_save(&config);
    else if (cmd == CONFIG_CMD_LOAD)
      // Load configuration from flash (undo any runtime changes).
      config_load(&config);
    else if (cmd == CONFIG_CMD_DEFAULTS)
      // Restore runtime configuration to compiled-in defaults.
      config_set_defaults(&config);
    else if (cmd == CONFIG_CMD_BOOTSEL)
      // Reboot device into bootsel mode.
      reset_usb_boot(0, 0);
  }
}
