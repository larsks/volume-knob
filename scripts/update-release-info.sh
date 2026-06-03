#!/bin/bash

vcs_ref=$(git rev-parse --short=10 HEAD || echo unknown)
vcs_date=$(git show -s --format=%ci HEAD)
vcs_version=$(git describe --tags)
version=${CZ_PRE_NEW_VERSION:-$vcs_version}

sed -i \
  -e "s/\.version =.*/.version = \"${version}\",/" \
  -e "s/\.vcs_ref =.*/.vcs_ref = \"${vcs_ref}\",/" \
  -e "s/\.vcs_date =.*/.vcs_date = \"${vcs_date}\",/" \
  version.c

sed -i \
  -e "s/version *string =.*/version string = \"${version}\"/" \
  -e "s/vcs_ref *string =.*/vcs_ref string = \"${vcs_ref}\"/" \
  -e "s/vcs_date *string =.*/vcs_date string = \"${vcs_date}\"/" \
  tools/version.go

git add version.c tools/version.go
