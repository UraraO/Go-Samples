package main

import (
	"fileserver/src/encryptor"
)

// curl -X POST -F "file=@./23.jpeg" "http://localhost:8080/upload/23.jpeg?algorithm=aes-gcm"
// curl -X GET --output "receivedfile.jpeg" "http://localhost:8080/download/23.jpeg?algorithm=aes-gcm"

func main() {
	// filespliter.FileSplitandMergeTest()
	encryptor.EncryptorTest()
}
