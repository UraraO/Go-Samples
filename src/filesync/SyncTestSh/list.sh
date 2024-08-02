#!/bin/bash

limit=$1
offset=$2

curl -v '127.0.0.1:8080/files?limit='$limit'&offset='$offset''
