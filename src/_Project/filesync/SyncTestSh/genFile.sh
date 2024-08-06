#!/bin/bash

size=$1
curDir=$(cd $(dirname $0); pwd)
filename=${curDir##*/}_$size
fileloc="./"$filename""
dd if=/dev/urandom of=$fileloc bs=$size count=1 &>/dev/null
