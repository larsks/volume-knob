#ifndef CONFIG_H
#define CONFIG_H

#define ENCODER_PIN_A 28
#define ENCODER_PIN_B 29

#define KEY_CW HID_USAGE_CONSUMER_VOLUME_INCREMENT
#define KEY_CCW HID_USAGE_CONSUMER_VOLUME_DECREMENT

// Number of detent clicks per key event
#define ENCODER_DIVIDER 256

#endif
