package hubwrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UsersList struct {
	Users []struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		LastActivity time.Time `json:"last_activity"`
		Email        string    `json:"email"`
	} `json:"users"`
}

func GetUsers() string {
	url := "https://api.hubstaff.com/v1/users"
	//var jsonStr = []byte(`{"email":"artem.potapov@mentalstack.com","password":"HubPass12~~"}`)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("app-token", "YQmK0Bo1-NU8u73q7M_prYrahzfcSfJJGSJbchBhP7k")
	req.Header.Set("auth-token", "zPM4B3F9PpOflINyv2dSYZSgKlks6IOg45xild8WfJ4")
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
		model := new(UsersList)
		json.NewDecoder(resp.Body).Decode(model)
		return model.Users[0].Email
	} else {
		return "-1"
	}
}
