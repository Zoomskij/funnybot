package hubwrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	ID                int         `json:"id"`
	Summary           string      `json:"summary"`
	Details           string      `json:"details"`
	RemoteID          string      `json:"remote_id"`
	RemoteAlternateID interface{} `json:"remote_alternate_id"`
	CompletedAt       interface{} `json:"completed_at"`
	Status            string      `json:"status"`
	ProjectID         int         `json:"project_id"`
}

type TaskList struct {
	Task []Task `json:"task"`
}
type TaskArr []Task

func GetTasks() string {
	url := "https://api.hubstaff.com/v1/tasks?offset=650"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("app-token", "YQmK0Bo1-NU8u73q7M_prYrahzfcSfJJGSJbchBhP7k")
	req.Header.Set("auth-token", "zPM4B3F9PpOflINyv2dSYZSgKlks6IOg45xild8WfJ4")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//var reply string
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)

		tasks := new(Task)
		json.NewDecoder(resp.Body).Decode(tasks)

		// for i := 0; i < len(tasks); i++ {
		// 	reply += tasks[i].RemoteAlternateID
		// }

		return "0"
	} else {
		return "-1"
	}
}
