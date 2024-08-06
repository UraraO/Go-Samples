package backup

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// GetPwd 给密码加密
func GetPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePwd 比对密码
func ComparePwd(pwd1 string, pwd2 string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(pwd1), []byte(pwd2))
	if err != nil {
		return false
	} else {
		return true
	}
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
