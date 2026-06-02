#include "version.h"

__attribute__((used, section(".version_info")))
const version_info_t version_info = {
    ._start = VERSION_START_MAKER,
    .version = "0.3.0",
    .vcs_ref = "3926197ffa",
    .vcs_date = "2026-06-02 18:14:51 -0400",
    ._end = VERSION_END_MAKER,
};
