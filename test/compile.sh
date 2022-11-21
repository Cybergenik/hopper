#!/bin/bash

CC="clang++-15"
#ASAN
CC+=" -fsanitize=address"
CC+=" -fno-omit-frame-pointer"
# Edge Coverage
CC+=" -fsanitize-coverage=edge,trace-pc-guard"
CC+=" -o target"

$CC -g $1



