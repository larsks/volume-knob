#ifndef PERSIST_H
#define PERSIST_H

#include <stdint.h>

// This is used to detect whether or not flash holds
// a valid configuration. It should be udpated whenever
// the config_t struct changes.
#define CONFIG_MAGIC 0x564B4332 // "VKC2"

typedef struct {
  uint32_t magic;
  uint8_t type_cw;
  uint8_t type_ccw;
  uint16_t key_cw;
  uint16_t key_ccw;
  uint16_t divider;
  uint32_t crc32;
} config_t;

void config_set_defaults(config_t *cfg);
void config_load(config_t *cfg);
void config_save(const config_t *cfg);

#endif
