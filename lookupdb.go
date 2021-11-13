package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
)

type NameModel struct {
	Name    string
	Address string
}

type Name struct {
	Name    string
	Address net.IP
}

func GetNames() ([]Name, error) {
	// read file
	data, err := ioutil.ReadFile("../htb_etc_hosts/hosts_all.txt")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var models []NameModel

	re := regexp.MustCompile(`(?m)^([^#]+?)\s+([^#]+?)\s*(?:#.*)$`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	models = make([]NameModel, 0)
	for _, match := range matches {
		model := NameModel{
			Name:    match[2],
			Address: match[1],
		}
		// fmt.Println("Matched", model)
		models = append(models, model)
	}
	// fmt.Println("Matched", len(models), "in total.")

	return To(models), nil
}

func To(models []NameModel) []Name {
	names := make([]Name, 0, len(models))
	for _, value := range models {
		ip := net.ParseIP(value.Address)
		// fmt.Println("Resolved", value.Address, "to", ip)
		names = append(names, Name{
			Name:    value.Name,
			Address: ip,
		})
	}

	// fmt.Println(names)
	return names
}
