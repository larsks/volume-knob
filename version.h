#ifndef _VERSION_H
#define _VERSION_H

#define VERSION_START_MAKER "_VERSION_START_"
#define VERSION_END_MAKER "_VERSION_END_"

// The _start and _end markers are used to identify version information
// in the output of the `strings` command.
typedef struct {
  char _start[16];
  char version[16];
  char vcs_ref[16];
  char vcs_date[32];
  char _end[16];
} version_info_t;

extern const version_info_t version_info;

#endif // __VERSION_H
