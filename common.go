package main

type Data struct {
	// Optional, also cache all data
	// Substitutes []structs.Substitute `json:"substitutes"`
	Meta Class `json:"meta"`
}

type Class struct {
	Class string `json:"class"`
	Date  string `json:"date"`
}

type Users []struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Class string `json:"class"`
}
