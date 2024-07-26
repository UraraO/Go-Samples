#! /bin/bash

#./fileserver &

sleep 1

curl -X POST -F "file=@./23.jpeg" "http://localhost:8080/upload/23.jpeg?algorithm=aes-gcm"
