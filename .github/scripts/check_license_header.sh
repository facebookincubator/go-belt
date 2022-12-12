#!/bin/sh
FILE_PATH="$1"; shift
exec grep -E '(Meta Platforms, Inc. and affiliates)|(Facebook, Inc(\.|,)? and its affiliates)|([0-9]{4}-present(\.|,)? Facebook)|([0-9]{4}(\.|,)? Facebook)' "$FILE_PATH" >/dev/null
