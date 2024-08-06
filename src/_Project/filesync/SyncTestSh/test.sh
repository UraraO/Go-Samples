#!/bin/bash

for((i=1;i<=3;i+=1));
do  
    size=$i"M"
    ./genFile.sh $size

    ./uploadFile.sh $size

    ./downloadFile.sh $size 8080
done