#!/bin/sh

if [ $# -eq 0 ]; then
  /bin/sh
else
  exec go-check "$@"
fi
