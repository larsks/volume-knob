#include "version.h"

__attribute__((used, section(".version_info")))
const version_info_t version_info = {
    ._start = VERSION_START_MAKER,
    .version = "0.2.0",
    .vcs_ref = "603d9d398b",
    .vcs_date = "2026-06-01 17:04:36 -0400",
    ._end = VERSION_END_MAKER,
};
