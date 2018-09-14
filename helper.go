package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
)

func getReceivers() Users {
	// Read in configuration file
	bytes, err := ioutil.ReadFile("users.json")

	if err != nil {
		log.Fatal("Failed to read users file: ", err)
	}

	var users Users
	if err := json.Unmarshal(bytes, &users); err != nil {
		log.Fatal(err)
	}

	return users
}
