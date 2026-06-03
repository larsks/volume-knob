#include "version.h"

__attribute__((used, section(".version_info")))
const version_info_t version_info = {
    ._start = VERSION_START_MAKER,
    .version = "1.1.0",
    .vcs_ref = "630ae44d10",
    .vcs_date = "2026-06-02 23:00:17 -0400",
    ._end = VERSION_END_MAKER,
};
