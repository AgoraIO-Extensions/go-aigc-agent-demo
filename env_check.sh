#!/usr/bin/env bash
set -e

output=$(gcc --version)

if [ -n "$output" ]; then
    echo "--ok-- gcc --version"
    echo "$output"
else
    echo "--fail-- ldd --version"
	exit 1
fi