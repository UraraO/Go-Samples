#!/bin/bash

size=$1
curDir=$(cd $(dirname $0); pwd)
filename=${curDir##*/}_$size
fileloc="./"$filename""

curl -X DELETE '127.0.0.1:8080/files/'$filename''