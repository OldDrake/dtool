package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func OutputJSON(data interface{}, filename string) error {
	jsonstr, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("JSON encoding error:", err)
		return err
	}
	if filename == "-" {
		fmt.Println(string(jsonstr))
	} else {
		err := ioutil.WriteFile(filename, jsonstr, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
