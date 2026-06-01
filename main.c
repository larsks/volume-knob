#include "bsp/board_api.h"
#include "hardware/gpio.h"
#include "pico/stdio.h"
#include "tusb.h"

#include "config.h"
#include "usb_hid.h"

// Quadrature state transition table.
// Index: (prev_AB << 2) | curr_AB
// Values: 0 = invalid/no change, 1 = CW, -1 = CCW
static const int8_t encoder_table[16] = {
    0, -1, 1, 0, 1, 0, 0, -1, -1, 0, 0, 1, 0, 1, -1, 0,
};

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

static bool send_consumer_key(uint16_t usage) {
  if (!tud_hid_ready())
    return false;
  tud_hid_report(REPORT_ID_CONSUMER_CONTROL, &usage, sizeof(usage));
  key_is_pressed = (usage != 0);
  return true;
}

static void encoder_task(void) {
  uint8_t curr = read_encoder_state();
  uint8_t index = (encoder_prev_state << 2) | curr;
  encoder_prev_state = curr;

  int8_t delta = encoder_table[index];
  if (delta == 0)
    return;

  encoder_accum += delta;

  if (encoder_accum >= ENCODER_DIVIDER) {
    if (send_consumer_key(KEY_CW))
      encoder_accum = 0;
  } else if (encoder_accum <= -ENCODER_DIVIDER) {
    if (send_consumer_key(KEY_CCW))
      encoder_accum = 0;
  }
}

static void release_task(void) {
  if (key_is_pressed && tud_hid_ready()) {
    send_consumer_key(0);
  }
}

int main(void) {
  board_init();
  tusb_init();
  stdio_init_all();
  encoder_init();

  while (1) {
    tud_task();
    encoder_task();
    release_task();
  }
}

void tud_mount_cb(void) {}
void tud_umount_cb(void) {}
void tud_suspend_cb(bool remote_wakeup_en) { (void)remote_wakeup_en; }
void tud_resume_cb(void) {}

uint16_t tud_hid_get_report_cb(uint8_t instance, uint8_t report_id,
                               hid_report_type_t report_type, uint8_t *buffer,
                               uint16_t reqlen) {
  (void)instance;
  (void)report_id;
  (void)report_type;
  (void)buffer;
  (void)reqlen;
  return 0;
}

void tud_hid_set_report_cb(uint8_t instance, uint8_t report_id,
                           hid_report_type_t report_type, uint8_t const *buffer,
                           uint16_t bufsize) {
  (void)instance;
  (void)report_id;
  (void)report_type;
  (void)buffer;
  (void)bufsize;
}
