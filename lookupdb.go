package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
)

type NameModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Name struct {
	Name    string
	Address net.IP
}

func GetNames() ([]Name, error) {
	// read file
	data, err := ioutil.ReadFile("./names.json")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	// json data
	var models []NameModel

	// unmarshall it
	err = json.Unmarshal(data, &models)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return To(models), nil

}

func To(models []NameModel) []Name {
	names := make([]Name, 0, len(models))
	for _, value := range models {
		names = append(names, Name{
			Name:    value.Name,
			Address: net.ParseIP(value.Address),
		})
	}
	return names
}
