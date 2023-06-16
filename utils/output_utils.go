package utils

import (
	"encoding/json"
	"fmt"
)

func OutputJSON(data map[string][]string) {
	jsonstr, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("JSON encoding error:", err)
		return
	}
	fmt.Println(string(jsonstr))
}
