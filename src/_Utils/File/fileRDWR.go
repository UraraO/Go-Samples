package backup

import (
	"fmt"
	"io"
	"os"
)

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadFile() []byte {
	f, err := os.OpenFile("file.txt", os.O_RDWR, os.ModeTemporary)
	if err != nil {
		fmt.Println("read file fail", err)
		return []byte{}
	}
	defer f.Close()

	fd, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return []byte{}
	}

	fmt.Println(string(fd))
	return fd
}

func WriteF(src []byte, fileName string) {
	err := os.WriteFile(fileName, src, 0666)
	if err != nil {
		fmt.Println("write fail")
	}
	fmt.Println("write success")
}

func fileRDWRtest() {
	fmt.Println(CheckFileExist("file.txt"))
	WriteF(ReadFile(), "file2.txt")
}
