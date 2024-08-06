/*===========
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:45:57
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-06 22:31:35
 FilePath: /Golang-Samples/src/_Library/JWT_authentication/JWT_test.go
 Description:

 JWT鉴权示例，包括Macke生成token，ParseToken解析token

 Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
===========*/

package backup

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Userinfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Macke 生成jwt 需要传入 用户名和密码
func Macke(user *Userinfo) (token string, err error) {
	claims := jwt.MapClaims{ // 创建一个自己的声明
		"name": user.Username,
		"pwd":  user.Password,
		"iss":  "lva",
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Second * 4).Unix(),
		"iat":  time.Now().Unix(),
	}
	then := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = then.SignedString([]byte("gettoken"))
	return
}

// secret 自己解析的秘钥
func secret() jwt.Keyfunc {
	// 按照这样的规则解析
	return func(t *jwt.Token) (interface{}, error) {
		return []byte("gettoken"), nil
	}
}

// 解析token
func ParseToken(token string) (user *Userinfo, err error) {
	user = &Userinfo{}
	ptoken, _ := jwt.Parse(token, secret())

	claim, ok := ptoken.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("解析错误")
		return
	}
	if !ptoken.Valid {
		err = errors.New("令牌错误！")
		return
	}
	// fmt.Println(claim)
	user.Username = claim["name"].(string) // 强行转换为string类型
	user.Password = claim["pwd"].(string)  // 强行转换为string类型
	return
}

func JWT_Test() {
	var use = Userinfo{"UraraO", "123456"}
	tkn, _ := Macke(&use)
	fmt.Println("_____", tkn)
	// time.Sleep(time.Second * 8)超过时间打印令牌错误
	user, err := ParseToken(tkn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user.Username, user.Password)
}
