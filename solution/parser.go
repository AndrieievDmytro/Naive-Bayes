package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readCsv(path string) ([][]string, error) {
	dataFile, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()
	if err == nil {
		buf, rerr := ioutil.ReadAll(dataFile)
		if rerr != nil {
			fmt.Println("Read CSV error: " + rerr.Error())
		}
		r := csv.NewReader(strings.NewReader(string(buf)))
		records, err := r.ReadAll()
		if err != nil {
			fmt.Println("Parse CSV error: " + err.Error())
		}
		return records, nil
	}
	return nil, err
}

func convertStrArrayToJson(records [][]string) string {
	jsonData := ""
	dimension := len(records[0]) - 1
	for _, record := range records {
		wrongStr := false
		if len(record) < dimension || len(record) > dimension+1 {
			fmt.Println("Wrong parameters count")
			wrongStr = true
		}
		// decisionAttribute := record[len(record)-1]
		attributes := record[:]
		stringArray := "["
		quotes := "\""
		for ind := range attributes {
			attributes[ind] = quotes + attributes[ind] + quotes
		}
		arrField := strings.Join(attributes, ",")
		stringArray += arrField
		stringArray += "]"
		if !wrongStr {
			jsonData += "{ \"Attributes\":" + stringArray + " },"
		}
	}
	jsonData = "[" + jsonData[:len(jsonData)-1] + "]"
	return jsonData
}

func (obj *Objects) readData(path string) {
	records, err := readCsv(path)
	jsonData := ""
	if err == nil {
		jsonData = convertStrArrayToJson(records)
		json.Unmarshal([]byte(jsonData), &obj.Obj)
	} else {
		fmt.Println("Read CSV error: " + err.Error())
	}
}
