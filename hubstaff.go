package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Auth struct {
	User struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		LastActivity time.Time `json:"last_activity"`
		AuthToken    string    `json:"auth_token"`
	} `json:"user"`
}

func hubauth() string {
	url := "https://api.hubstaff.com/v1/auth"
	var jsonStr = []byte(`{"email":"artem.potapov@mentalstack.com","password":"HubPass12~~"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("app-token", "YQmK0Bo1-NU8u73q7M_prYrahzfcSfJJGSJbchBhP7k")
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
		model := new(Auth)
		json.NewDecoder(resp.Body).Decode(model)
		return model.User.AuthToken
	} else {
		return "-1"
	}
}
