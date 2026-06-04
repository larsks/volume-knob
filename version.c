#include "version.h"

__attribute__((used, section(".version_info")))
const version_info_t version_info = {
    ._start = VERSION_START_MAKER,
    .version = "1.2.0",
    .vcs_ref = "c1a7b639e7",
    .vcs_date = "2026-06-03 22:59:03 -0400",
    ._end = VERSION_END_MAKER,
};
