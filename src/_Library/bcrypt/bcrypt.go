/*
  - @Author: chaidaxuan chaidaxuan@wps.cn
  - @Date: 2024-07-26 16:58:52
 * @LastEditors: chaidaxuan chaidaxuan@wps.cn
 * @LastEditTime: 2024-08-06 09:47:15
 * @FilePath: /Golang-Samples/src/_Library/bcrypt/bcrypt.go
  - @Description:

"golang.org/x/crypto/bcrypt"

本文件展示bcrypt库的使用示例
bcrypt库实现了Provos and Mazières's bcrypt adaptive hashing algorithm
用于从一串字符创建密码hash，并用原始密码与hash串对比

  - Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
*/

package backup

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GetPwd 给密码明文加密
func GetPwd(pwdPlaintext string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwdPlaintext), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePwd 比对密码
func ComparePwd(pwd1 string, pwd2 string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(pwd1), []byte(pwd2))
	return err == nil
}

func BcryptTest() {
	password := "Urara Password"
	pwd1, err := GetPwd(password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("password =", password)
	fmt.Println("pwd1 =", pwd1)
	pswValid := ComparePwd(pwd1, password)
	if pswValid {
		fmt.Println("password validated pass~")
	} else {
		fmt.Println("password validated error!")
	}
}
