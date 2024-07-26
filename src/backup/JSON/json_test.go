package backup

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type UserSend2 struct {
	Username string
	Password string
	Somes    some
}

type some struct {
	Somestring string
}

func JsonTest() {
	key1 := base64.StdEncoding.EncodeToString([]byte("thisisplainkey"))
	key2 := base64.StdEncoding.EncodeToString([]byte(key1))
	fmt.Println(key1)
	fmt.Println(key2)
	de1, _ := base64.StdEncoding.DecodeString(key2)
	de2, _ := base64.StdEncoding.DecodeString(key1)
	fmt.Println(string(de1))
	fmt.Println(string(de2))

	user := UserSend2{
		Username: "Urara_json_send1",
		Password: "password",
		Somes: some{
			Somestring: "astring",
		},
	}
	body, _ := json.Marshal(user)
	fmt.Println("body:", string(body))

	var ret UserSend2
	err := json.Unmarshal(body, &ret)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(ret)
	fmt.Println(ret.Somes.Somestring)
	de3, _ := base64.StdEncoding.DecodeString("MTIzNDU2NzgxMjM0NTY3OA==")
	fmt.Println(string(de3))
}
