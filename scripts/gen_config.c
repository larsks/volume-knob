#include <stddef.h>
#include <stdio.h>

#include "persist.h"
#include "usb_hid.h"

#define FIELD(name) \
  printf("field %s %zu %zu\n", #name, offsetof(config_t, name), sizeof(((config_t *)0)->name))

#define CONST_HEX(name) printf("const %s 0x%X\n", #name, name)
#define CONST_DEC(name) printf("const %s %d\n", #name, name)

int main(void) {
  printf("report_size %zu\n", CONFIG_REPORT_SIZE);

  FIELD(magic);
  FIELD(type_cw);
  FIELD(type_ccw);
  FIELD(key_cw);
  FIELD(key_ccw);
  FIELD(divider);
  FIELD(mod_cw);
  FIELD(mod_ccw);
  FIELD(crc32);

  CONST_HEX(CONFIG_MAGIC);
  CONST_DEC(REPORT_ID_CONFIG);
  CONST_DEC(REPORT_ID_COMMAND);
  CONST_DEC(CONFIG_CMD_SAVE);
  CONST_DEC(CONFIG_CMD_LOAD);
  CONST_DEC(CONFIG_CMD_DEFAULTS);
  CONST_DEC(CONFIG_CMD_BOOTSEL);
  CONST_DEC(KEY_TYPE_CONSUMER);
  CONST_DEC(KEY_TYPE_KEYBOARD);
  CONST_DEC(COMMAND_REPORT_SIZE);

  return 0;
}
