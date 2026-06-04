#ifndef PERSIST_H
#define PERSIST_H

#include <stdint.h>

// Convert a four character sequence into a 4-byte integer constant.
#define FOURCC(a, b, c, d)                                                     \
  ((uint32_t)(a) << 24 | (uint32_t)(b) << 16 | (uint32_t)(c) << 8 |            \
   (uint32_t)(d))

// This is used to detect whether or not flash holds
// a valid configuration. It should be updated whenever
// the config_t struct changes.
#define CONFIG_MAGIC FOURCC('V', 'K', 'C', '3')

typedef struct __attribute__((packed)) {
  uint32_t magic;
  uint8_t type_cw;
  uint8_t type_ccw;
  uint16_t key_cw;
  uint16_t key_ccw;
  uint16_t divider;
  uint8_t mod_cw;
  uint8_t mod_ccw;
  uint32_t crc32;
} config_t;

void config_set_defaults(config_t *cfg);
void config_load(config_t *cfg);
void config_save(const config_t *cfg);

#endif
