package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	Port        = "7474"
	Ip          = "10.0.0.213"
	FaceActions = []Action{}
)

type FaceEnum string

const (
	Thug    FaceEnum = "thug"
	Ironman FaceEnum = "ironman"
	Batman  FaceEnum = "batman"
	Clown   FaceEnum = "clown"
	Eyes    FaceEnum = "eyes"
	Frog    FaceEnum = "frog"
	Likes   FaceEnum = "Likes" // Likes action to map streambot likes action
	Cry     FaceEnum = "Cry"
)

//	{
//		"id": "aedb1356-a9c7-4a0c-8f8a-030d7f24612d",
//		"name": "Likes",
//		"group": "",
//		"enabled": true,
//		"subactions_count": 6
//	  },
type Response struct {
	Count   int      `json:"count"`
	Actions []Action `json:"actions"`
}
type Action struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

func GetActionList() {
	resp, err := http.Get("http://" + Ip + ":" + Port + "/GetActions")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	fmt.Println(sb)
	in := []byte(string(sb))

	var iot Response
	temp := json.Unmarshal(in, &iot)
	// fmt.Println("temp:", temp)
	if temp != nil {
		panic(err)
	}
	// fmt.Println("iot:", iot)
	for _, action := range iot.Actions {
		if action.Group == "face" {
			FaceActions = append(FaceActions, action)
		}
	}
	fmt.Println("FaceActions:", FaceActions)

}

func DoAction(action Action) {
	body := map[string]interface{}{
		"action": map[string]interface{}{
			"id":   action.ID,
			"name": action.Name,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	resp, err := http.Post("http://"+Ip+":"+Port+"/DoAction", "application/json", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		log.Fatalln("Error sending POST request:", err)
	}
	defer resp.Body.Close()

}
func GetAction(name string) Action {
	output := Action{}
	for _, action := range FaceActions {
		if action.Name == name {
			output = action
		}
	}
	fmt.Println("output:", output)
	return output
}
