package hubwrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthModel struct {
	User struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		LastActivity time.Time `json:"last_activity"`
		AuthToken    string    `json:"auth_token"`
	} `json:"user"`
}

func Auth() string {
	url := "https://api.hubstaff.com/v1/auth"
	var jsonStr = []byte(`{"email":"EMAIL","password":"PASSWORD"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("app-token", "TOKEN")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		model := new(AuthModel)
		json.NewDecoder(resp.Body).Decode(model)
		return model.User.AuthToken
	} else {
		return "-1"
	}
}
