package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserSend struct {
	Username string
	Password string
}

func UserSendTest() {
	user := UserSend{
		Username: "Urara_json_send1",
		Password: "password",
	}
	body, _ := json.Marshal(user)
	resp, err := http.Post("http://127.0.0.1:8080/api/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)

	} else {
		fmt.Println("Get failed with error: ", resp.Status)
	}
}
