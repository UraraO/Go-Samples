#!/bin/bash

size=$1
port=$2

curDir=$(cd $(dirname $0); pwd)
filename=${curDir##*/}_$size
originFileloc="./"$filename""
downloadFileLoc="./"$filename"_download_"$port""

curl -o $downloadFileLoc '127.0.0.1:'$port'/files/'$filename''

if ! cmp -s $downloadFileLoc $originFileloc ; then
    echo "upload and download not same, size: "$size
    exit 1
fi