#ifndef CONFIG_H
#define CONFIG_H

// These values define the GPIOs to which the encoder outputs are connected
#define ENCODER_PIN_A 28
#define ENCODER_PIN_B 29

// Key sent on clockwise motion
#define KEY_CW HID_USAGE_CONSUMER_VOLUME_INCREMENT

// Key sent on counter-clockwise motion
#define KEY_CCW HID_USAGE_CONSUMER_VOLUME_DECREMENT

// Number of transitions before sending a key event
#define ENCODER_DIVIDER 256

#endif
