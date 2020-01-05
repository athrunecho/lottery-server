package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type EmployeeFormat struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PrizeFormat struct {
	Name string `json:"name"`
	Num  int    `json:"num"`
}

func loadEmployeeConfig(file string, e *[]EmployeeFormat) {

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = json.Unmarshal(buf, e)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func loadPrizeConfig(file string, p *[]PrizeFormat) {

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = json.Unmarshal(buf, p)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func getConfig() ([]EmployeeFormat, []PrizeFormat) {
	ePath := path.Join("./", "employees.json")
	pPath := path.Join("./", "prizes.json")

	loadEmployeeConfig(ePath, &employeeArray)
	loadPrizeConfig(pPath, &prizeArray)
	return employeeArray, prizeArray
}
