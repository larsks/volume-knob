#include "bsp/board_api.h"
#include "hardware/gpio.h"
#include "pico/stdio.h"
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
static bool key_is_pressed;

static uint8_t read_encoder_state(void) {
  uint8_t a = gpio_get(ENCODER_PIN_A) ? 1 : 0;
  uint8_t b = gpio_get(ENCODER_PIN_B) ? 1 : 0;
  return (a << 1) | b;
}

static void encoder_init(void) {
  gpio_init(ENCODER_PIN_A);
  gpio_set_dir(ENCODER_PIN_A, GPIO_IN);
  gpio_pull_up(ENCODER_PIN_A);

  gpio_init(ENCODER_PIN_B);
  gpio_set_dir(ENCODER_PIN_B, GPIO_IN);
  gpio_pull_up(ENCODER_PIN_B);

  encoder_prev_state = read_encoder_state();
  encoder_accum = 0;
  key_is_pressed = false;
}

// Generate a key event if the USB HID endpoint is ready to accept a new
// report. Return true if the endpoint was ready, false if it was not and we
// did not send a key.
static bool send_consumer_key(uint16_t usage) {
  if (!tud_hid_ready())
    return false;
  tud_hid_report(REPORT_ID_CONSUMER_CONTROL, &usage, sizeof(usage));
  key_is_pressed = (usage != 0);
  return true;
}

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
    if (send_consumer_key(config.key_cw))
      encoder_accum = 0;
  } else if (encoder_accum <= -config.divider) {
    if (send_consumer_key(config.key_ccw))
      encoder_accum = 0;
  }
}

static void release_task(void) {
  if (key_is_pressed)
    send_consumer_key(0);
}

int main(void) {
  board_init();
  tusb_init();
  stdio_init_all();
  config_load(&config);
  encoder_init();

  while (1) {
    tud_task();
    encoder_task();
    release_task();
  }
}

void tud_suspend_cb(bool remote_wakeup_en) { (void)remote_wakeup_en; }

// In order for everything to build properly, we must define these functions
// even if we're not implementing them.
void tud_mount_cb(void) {}
void tud_umount_cb(void) {}
void tud_resume_cb(void) {}

uint16_t tud_hid_get_report_cb(uint8_t instance, uint8_t report_id,
                               hid_report_type_t report_type, uint8_t *buffer,
                               uint16_t reqlen) {
  (void)report_type;
  (void)reqlen;

  if (instance == HID_INSTANCE_CONFIG && report_id == REPORT_ID_CONFIG) {
    memcpy(buffer, &config.key_cw, CONFIG_REPORT_SIZE);
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
    memcpy(&config.key_cw, buffer, CONFIG_REPORT_SIZE);
  } else if (report_id == REPORT_ID_COMMAND && bufsize >= COMMAND_REPORT_SIZE) {
    uint8_t cmd = buffer[0];
    if (cmd == CONFIG_CMD_SAVE)
      config_save(&config);
    else if (cmd == CONFIG_CMD_LOAD)
      config_load(&config);
    else if (cmd == CONFIG_CMD_DEFAULTS)
      config_set_defaults(&config);
  }
}
