#!/usr/bin/env bash
set -e

output=$(ldd --version)

if [ -n "$output" ]; then
    echo "--ok-- ldd --version"
    echo "$output"
else
    echo "--fail-- ldd --version"
	exit 1
fi

output=$(ldconfig -p | grep libpthread)

if [ -n "$output" ]; then
    echo "--ok-- ldconfig -p | grep libpthread"
    echo "$output"
else
    echo "--fail-- ldconfig -p | grep libpthread"
	exit 1
fi

output=$(openssl version)

if [ -n "$output" ]; then
    echo "--ok-- openssl version"
    echo "$output"
else
    echo "--fail-- openssl version"
	exit 1
fi

output=$(dpkg -l | grep libssl)

if [ -n "$output" ]; then
    echo "--ok-- dpkg -l | grep libssl"
    echo "$output"
else
    echo "--fail-- dpkg -l | grep libssl"
    exit 1
fi

output=$(dpkg -l | grep ca-certificates)

if [ -n "$output" ]; then
    echo "--ok-- dpkg -l | grep ca-certificates"
    echo "$output"--
else
    echo "--fail-- dpkg -l | grep libssl"
    exit 1
fi