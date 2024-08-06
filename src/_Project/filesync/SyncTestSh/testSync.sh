#!/bin/bash
# size=$1

# ./genFile.sh $size

# ./uploadFile.sh $size

# ./downloadFile.sh $size 8080

# ./downloadFile.sh $size 8082

# for((i=1;i<=3;i+=1));
# do  
#     size=$i"M"
#     ./genFile.sh $size

#     ./uploadFile.sh $size

#     ./downloadFile.sh $size 8080
# done

for((i=1;i<=3;i+=1));
do  
    size=$i"M"

    ./downloadFile.sh $size 8082
done

