#ifndef CONFIG_H
#define CONFIG_H

#include "tusb.h"

#define ENCODER_PIN_A 28
#define ENCODER_PIN_B 29

#define KEY_CW HID_USAGE_CONSUMER_VOLUME_INCREMENT
#define KEY_CCW HID_USAGE_CONSUMER_VOLUME_DECREMENT

// State transitions per detent click (4 for most quadrature encoders)
#define ENCODER_STEPS_PER_DETENT 4

// Number of detent clicks per key event
#define ENCODER_DIVIDER 64

#endif
