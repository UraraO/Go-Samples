#! /bin/bash

#./fileserver &

sleep 1

curl -X GET --output "receivedfile.jpeg" "http://localhost:8080/download/23.jpeg?algorithm=aes-gcm"

echo "remember to remove the receivedfile.jpeg"