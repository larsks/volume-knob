#ifndef CONFIG_H
#define CONFIG_H

#include "tusb.h"

#define ENCODER_PIN_A 28
#define ENCODER_PIN_B 29

#define KEY_CW HID_USAGE_CONSUMER_VOLUME_INCREMENT
#define KEY_CCW HID_USAGE_CONSUMER_VOLUME_DECREMENT

// Number of encoder state transitions per key event.
// Most encoders produce 4 transitions per detent; set to 4 for one key per
// click.
#define ENCODER_DIVIDER 4

#endif
