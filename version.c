#include "version.h"

__attribute__((used, section(".version_info")))
const version_info_t version_info = {
    ._start = VERSION_START_MAKER,
    .version = "0.2.0",
    .vcs_ref = "6210eacf57",
    .vcs_date = "2026-06-01 16:26:08 -0400",
    ._end = VERSION_END_MAKER,
};
