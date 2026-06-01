#include "persist.h"
#include "config.h"
#include "usb_hid.h"

#include "hardware/flash.h"
#include "hardware/sync.h"
#include "pico/stdlib.h"

#include <string.h>

#define CONFIG_FLASH_OFFSET (PICO_FLASH_SIZE_BYTES - FLASH_SECTOR_SIZE)

static const config_t *flash_config =
    (const config_t *)(XIP_BASE + CONFIG_FLASH_OFFSET);

static uint32_t crc32(const void *data, size_t len) {
  const uint8_t *p = data;
  uint32_t crc = 0xFFFFFFFF;
  for (size_t i = 0; i < len; i++) {
    crc ^= p[i];
    for (int j = 0; j < 8; j++) {
      crc = (crc >> 1) ^ (0xEDB88320 & -(crc & 1));
    }
  }
  return ~crc;
}

void config_set_defaults(config_t *cfg) {
  cfg->magic = CONFIG_MAGIC;
  cfg->type_cw = KEY_TYPE_CONSUMER;
  cfg->type_ccw = KEY_TYPE_CONSUMER;
  cfg->key_cw = DEFAULT_KEY_CW;
  cfg->key_ccw = DEFAULT_KEY_CCW;
  cfg->divider = DEFAULT_ENCODER_DIVIDER;
  cfg->crc32 = 0;
}

void config_load(config_t *cfg) {
  if (flash_config->magic != CONFIG_MAGIC) {
    config_set_defaults(cfg);
    return;
  }
  uint32_t expected = crc32(flash_config, offsetof(config_t, crc32));
  if (flash_config->crc32 != expected) {
    config_set_defaults(cfg);
    return;
  }
  memcpy(cfg, flash_config, sizeof(config_t));
}

static void __no_inline_not_in_flash_func(do_flash_write)(const config_t *cfg) {
  uint32_t ints = save_and_disable_interrupts();
  flash_range_erase(CONFIG_FLASH_OFFSET, FLASH_SECTOR_SIZE);
  flash_range_program(CONFIG_FLASH_OFFSET, (const uint8_t *)cfg,
                      sizeof(config_t));
  restore_interrupts(ints);
}

void config_save(const config_t *cfg) {
  config_t tmp = *cfg;
  tmp.magic = CONFIG_MAGIC;
  tmp.crc32 = crc32(&tmp, offsetof(config_t, crc32));
  do_flash_write(&tmp);
}
