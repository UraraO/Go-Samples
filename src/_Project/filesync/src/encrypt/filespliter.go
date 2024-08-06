package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

const chunkSize = 40000000

// const chunkSize = 50000000
var (
	action  string
	infile  string
	outfile string
)

// 最基础的切分文件函数
// 传入需要切分的源文件名
// 输出分块数量
func Split(infile string) int {
	if infile == "" {
		panic("请输入正确的文件名")
	}

	fileInfo, err := os.Stat(infile)
	if err != nil {
		if os.IsNotExist(err) {
			panic("文件不存在")
		}
		panic(err)
	}

	num := math.Ceil(float64(fileInfo.Size()) / chunkSize)

	fi, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Printf("要拆分成%.0f份\n", num)
	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(num); i++ {
		fi.Seek((i-1)*chunkSize, 0)
		if len(b) > int(fileInfo.Size()-(i-1)*chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)
		ofile := fmt.Sprintf("./%d.part", i)
		fmt.Printf("生成%s\n", ofile)
		f, err := os.OpenFile(ofile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f.Write(b)
		f.Close()
	}
	fi.Close()
	fmt.Println("拆分完成")
	return int(num)
}

// 最基础的合并文件函数
// 传入输出文件名
// 输出分块数量
func Merge(outfile string) int {
	fii, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	part_list, err := filepath.Glob("./*.part")
	if err != nil {
		panic(err)
	}
	fmt.Printf("要把%v份合并成一个文件%s\n", part_list, outfile)
	i := 0
	for _, v := range part_list {
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		fii.Write(b)
		f.Close()
		i++
		fmt.Printf("合并%d个\n", i)
	}
	fii.Close()
	fmt.Println("合并成功")
	return len(part_list)
}

// 切分加密文件函数
// 传入需要切分的源文件名
// 输出分块数量
func SplitEnc(infile string) int {
	if infile == "" {
		panic("请输入正确的文件名")
	}

	fileInfo, err := os.Stat(infile)
	if err != nil {
		if os.IsNotExist(err) {
			panic("文件不存在")
		}
		panic(err)
	}

	num := math.Ceil(float64(fileInfo.Size()) / chunkSize)

	fi, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Printf("要拆分成%.0f份\n", num)
	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(num); i++ {
		fi.Seek((i-1)*chunkSize, 0)
		if len(b) > int(fileInfo.Size()-(i-1)*chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)
		ofile := fmt.Sprintf("./%d.part.enc", i)
		fmt.Printf("生成%s\n", ofile)
		f, err := os.OpenFile(ofile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f.Write(b)
		f.Close()
	}
	fi.Close()
	os.Remove(fi.Name()) ////////////////////////////////////////
	fmt.Println("拆分完成")
	return int(num)
}

// 合并加密文件函数
// 传入输出文件名
// 输出分块数量
func MergeEnc(outfile string) int {
	fii, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	part_list, err := filepath.Glob("./*.part.enc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("要把%v份合并成一个文件%s\n", part_list, outfile)
	i := 0
	for _, v := range part_list {
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		fii.Write(b)
		f.Close()
		os.Remove(v) ////////////////////////////////////////
		i++
		fmt.Printf("合并%d个\n", i)
	}
	fii.Close()
	fmt.Println("合并成功")
	return len(part_list)
}

// 切分文件并加密
// 传入需要切分的源文件名和加密器
// 输出分块数量
// 保存的小文件是切分后再加密的文件
func SplitSourceToEnc(infile string, encryptor Encryptor) int {
	if infile == "" {
		panic("请输入正确的文件名")
	}

	fileInfo, err := os.Stat(infile)
	if err != nil {
		if os.IsNotExist(err) {
			panic("文件不存在")
		}
		panic(err)
	}

	num := math.Ceil(float64(fileInfo.Size()) / chunkSize)

	fi, err := os.OpenFile(infile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	fmt.Printf("要拆分成%.0f份\n", num)
	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(num); i++ {
		fi.Seek((i-1)*chunkSize, 0)
		if len(b) > int(fileInfo.Size()-(i-1)*chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)
		ofile := fmt.Sprintf("./%d.part.enc", i)
		fmt.Printf("生成%s\n", ofile)
		f, err := os.OpenFile(ofile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}

		// 加密文件内容
		encryptedPart, err := encryptor.Encrypt(b)
		if err != nil {
			panic(err)
		}

		f.Write([]byte(encryptedPart))
		f.Close()
	}
	fi.Close()
	os.Remove(infile) ////////////////////////////////////////
	fmt.Println("拆分完成")
	return int(num)
}

// 合并加密小文件并解密成源文件
// 传入输出文件名和解密器
// 输出分块数量
// 保存的文件是合并加密小文件后再解密的源文件
func MergeEncToSource(outfile string, encryptor Encryptor) int {
	fii, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	part_list, err := filepath.Glob("./*.part.enc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("要把%v份合并成一个文件%s\n", part_list, outfile)
	i := 0
	for _, v := range part_list {
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return 0
		}

		// 解密文件内容
		decryptedPart, err := encryptor.Decrypt(string(b))
		if err != nil {
			panic(err)
		}

		fii.Write(decryptedPart)
		f.Close()
		i++
		fmt.Printf("合并%d个\n", i)
	}
	fii.Close()
	fmt.Println("合并成功")
	return len(part_list)
}

// 合并加密小文件并解密成源文件，并通过http连接发送
// 传入输出文件名和解密器
// 输出分块数量
// 保存的文件是合并加密小文件后再解密的源文件
// 最终会删除源文件，防止明文留档
func MergeEncToSourceWriteHttp(outfile string, encryptor Encryptor, w http.ResponseWriter) int {
	fii, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	part_list, err := filepath.Glob("./*.part.enc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("要把%v份合并成一个文件%s\n", part_list, outfile)
	i := 0
	for _, v := range part_list {
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return 0
		}

		// 解密文件内容
		decryptedPart, err := encryptor.Decrypt(string(b))
		if err != nil {
			panic(err)
		}

		// fii.Write(decryptedPart) //////////////////////////////////////////////
		os.Remove(v) ////////////////////////////////////////
		w.Write(decryptedPart)
		f.Close()
		i++
		fmt.Printf("合并%d个\n", i)
	}
	fii.Close()
	os.Remove(fii.Name()) ////////////////////////////////////////
	fmt.Println("合并成功")
	return len(part_list)
}

// Encryptor 定义加密和解密的接口
type Encryptor interface {
	Encrypt(data []byte) (string, error)
	Decrypt(cryptoText string) ([]byte, error)
}

// AESGCMEncryptor 使用 AES-GCM 加密
type AESGCMEncryptor struct {
	key []byte
}

func NewAESGCMEncryptor(key []byte) *AESGCMEncryptor {
	return &AESGCMEncryptor{key: key}
}

func (e *AESGCMEncryptor) Encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *AESGCMEncryptor) Decrypt(cryptoText string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// AESCBCEncryptor 使用 AES-CBC 加密
type AESCBCEncryptor struct {
	key []byte
}

func NewAESCBCEncryptor(key []byte) *AESCBCEncryptor {
	return &AESCBCEncryptor{key: key}
}

func (e *AESCBCEncryptor) Encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	data = pkcs7Padding(data, block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *AESCBCEncryptor) Decrypt(cryptoText string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return pkcs7Unpadding(ciphertext)
}

// PKCS#7 padding and unpadding
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	padding := data[length-1]
	return data[:length-int(padding)], nil
}

func FileSplitandMergeTest() {
	flag.StringVar(&action, "a", "split", "请输入用途：split/merge 默认是split")
	flag.StringVar(&infile, "f", "", "请输入文件名")
	flag.StringVar(&outfile, "o", "azhang.mp4", "请输入要合并的文件名")
	flag.Parse()
	if action == "split" {
		Split(infile)
	} else if action == "merge" {
		Merge(outfile)
	} else {
		panic("-a只能输入split/merge")
	}
}
