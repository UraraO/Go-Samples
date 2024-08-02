package encrypt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var key = []byte("mysecretencryptionkey1234567890a")

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从请求参数中获取加密算法
	algorithm := r.URL.Query().Get("algorithm")
	var encryptor Encryptor
	switch algorithm {
	case "aes-gcm":
		encryptor = NewAESGCMEncryptor(key)
	case "aes-cbc":
		encryptor = NewAESCBCEncryptor(key)
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	// 读取文件
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 读取文件内容
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file data", http.StatusInternalServerError)
		return
	}

	// 加密文件内容
	encryptedData, err := encryptor.Encrypt(fileData)
	if err != nil {
		http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
		return
	}

	// 存储文件
	filePath := strings.TrimPrefix(r.URL.Path, "/upload/")
	if strings.HasSuffix(filePath, "/") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile("./"+filePath, []byte(encryptedData), 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File uploaded successfully"))
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从请求参数中获取加密算法
	algorithm := r.URL.Query().Get("algorithm")
	var encryptor Encryptor
	switch algorithm {
	case "aes-gcm":
		encryptor = NewAESGCMEncryptor(key)
	case "aes-cbc":
		encryptor = NewAESCBCEncryptor(key)
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/download/")
	if strings.HasSuffix(filePath, "/") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 读取文件
	encryptedData, err := ioutil.ReadFile("./" + filePath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// 解密文件内容
	decryptedData, err := encryptor.Decrypt(string(encryptedData))
	if err != nil {
		http.Error(w, "Failed to decrypt file", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(decryptedData)
}

func SplitUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从请求参数中获取加密算法
	algorithm := r.URL.Query().Get("algorithm")
	var encryptor Encryptor
	switch algorithm {
	case "aes-gcm":
		encryptor = NewAESGCMEncryptor(key)
	case "aes-cbc":
		encryptor = NewAESCBCEncryptor(key)
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	// 读取文件
	file, fileheader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 读取文件内容
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file data", http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile("./"+fileheader.Filename, fileData, 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// 将源文件切分为加密小文件
	SplitSourceToEnc(fileheader.Filename, encryptor)
	// 加密小文件合并为加密大文件，再删除中间过程的加密小文件
	MergeEnc(fileheader.Filename + ".enc")
	// // 切分Split源文件，修改下方读取和加密，改为批量，获取批量加密文件
	// // 修改存储，将批量加密后的文件Merge为完整加密文件
	// // fileheader.Filename
	// filenum := Split(fileheader.Filename)

	// for i := 0; i < filenum; i++ {

	// }

	// // 读取文件内容
	// fileData, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	http.Error(w, "Failed to read file data", http.StatusInternalServerError)
	// 	return
	// }

	// // 加密文件内容
	// encryptedData, err := encryptor.Encrypt(fileData)
	// if err != nil {
	// 	http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
	// 	return
	// }

	// // 存储文件
	// filePath := strings.TrimPrefix(r.URL.Path, "/upload/")
	// if strings.HasSuffix(filePath, "/") {
	// 	http.Error(w, "Invalid file path", http.StatusBadRequest)
	// 	return
	// }

	// err = ioutil.WriteFile("./"+filePath, []byte(encryptedData), 0644)
	// if err != nil {
	// 	http.Error(w, "Failed to save file", http.StatusInternalServerError)
	// 	return
	// }

	w.Write([]byte("File uploaded successfully"))
}

func SplitDownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从请求参数中获取加密算法
	algorithm := r.URL.Query().Get("algorithm")
	var encryptor Encryptor
	switch algorithm {
	case "aes-gcm":
		encryptor = NewAESGCMEncryptor(key)
	case "aes-cbc":
		encryptor = NewAESCBCEncryptor(key)
	default:
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/download/")
	if strings.HasSuffix(filePath, "/") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 加密文件切分为加密小文件
	SplitEnc(filePath + ".enc")

	// // 切分Split加密文件，修改下方读取和解密，改为批量，获取批量解密文件
	// // 修改存储，将批量解密后的文件Merge为完整源文件

	// // 读取文件
	// encryptedData, err := ioutil.ReadFile("./" + filePath)
	// if err != nil {
	// 	http.Error(w, "Failed to read file", http.StatusInternalServerError)
	// 	return
	// }

	// // 解密文件内容
	// decryptedData, err := encryptor.Decrypt(string(encryptedData))
	// if err != nil {
	// 	http.Error(w, "Failed to decrypt file", http.StatusInternalServerError)
	// 	return
	// }

	// 设置响应头
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	// 加密小文件合并为加密大文件，再解密并发送，中间过程的加密小文件和源文件不保留，直接通过http发送
	MergeEncToSourceWriteHttp(filePath, encryptor, w)
	// w.Write(decryptedData)
}

func EncryptorTest() {
	// http.HandleFunc("/upload/", UploadHandler)
	// http.HandleFunc("/download/", DownloadHandler)
	http.HandleFunc("/upload/", SplitUploadHandler)
	http.HandleFunc("/download/", SplitDownloadHandler)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
