#!/bin/bash

vcs_ref=$(git rev-parse --short=10 HEAD || echo unknown)
vcs_date=$(git show -s --format=%ci HEAD)

sed -i \
  -e "s/\.version =.*/.version = \"${CZ_PRE_NEW_VERSION}\",/" \
  -e "s/\.vcs_ref =.*/.vcs_ref = \"${vcs_ref}\",/" \
  -e "s/\.vcs_date =.*/.vcs_date = \"${vcs_date}\",/" \
  version.c
git add version.c
