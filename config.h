#ifndef CONFIG_H
#define CONFIG_H

#include "tusb.h"

// These values define the GPIOs to which the encoder outputs are connected
#define ENCODER_PIN_A 28
#define ENCODER_PIN_B 29

// Default key sent on clockwise motion
#define DEFAULT_KEY_CW HID_USAGE_CONSUMER_VOLUME_INCREMENT

// Default key sent on counter-clockwise motion
#define DEFAULT_KEY_CCW HID_USAGE_CONSUMER_VOLUME_DECREMENT

// Default number of transitions before sending a key event
#define DEFAULT_ENCODER_DIVIDER 256

#endif
