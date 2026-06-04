#include "version.h"

__attribute__((used, section(".version_info")))
const version_info_t version_info = {
    ._start = VERSION_START_MAKER,
    .version = "1.2.1",
    .vcs_ref = "93ed8fea8c",
    .vcs_date = "2026-06-04 08:41:55 -0400",
    ._end = VERSION_END_MAKER,
};
