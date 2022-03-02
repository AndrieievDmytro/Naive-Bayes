package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isContains(array [][]string, itterCount int, attr string) bool {
	answ := false
	for i := 0; i < len(array[itterCount]); i++ {
		if attr == array[itterCount][i] {
			answ = true
		}
	}
	return answ
}

func (obj *Objects) getUniqueAttributes() [][]string {
	arrays := make([][]string, len(obj.Obj[0].Attributes))
	for _, attr := range obj.Obj {
		atrs := attr.Attributes
		for i := 0; i < len(atrs); i++ {
			atr := atrs[i]
			if !isContains(arrays, i, atr) { //check for unique attribute in each column of data set
				arrays[i] = append(arrays[i], atr)
			}
		}
	}
	return arrays
}

func (obj *Objects) classify(attrs []string) string {
	attributeCount := make(map[string][]int, len(obj.Obj[0].Attributes)-1) // Repetition of attributes regarding to decision attribute
	occurenseCount := make(map[string]int, len(obj.Obj[0].Attributes)-1)   // Decision attribute occurence within data set
	attrLen := len(obj.Obj[0].Attributes) - 1                              //Attribute lenght without decision attribute
	for ind := range obj.Obj {
		attrcount := make([]int, attrLen)
		for i := 0; i < len(attrs); i++ {
			if obj.Obj[ind].Attributes[i] == attrs[i] {
				attrcount[i] += 1
			} else {
				attrcount[i] += 0
			}
		}
		decisionAttribute := obj.Obj[ind].Attributes[attrLen] // Decision attribute
		if attributeCount[decisionAttribute] == nil {
			attributeCount[decisionAttribute] = attrcount
		} else {
			prevAttributeCount := attributeCount[decisionAttribute]
			newAttributeCount := make([]int, attrLen)
			for j := 0; j < len(attrcount); j++ {
				newAttributeCount[j] += attrcount[j] + prevAttributeCount[j]
			}
			attributeCount[decisionAttribute] = newAttributeCount
		}
		occurenseCount[decisionAttribute] += 1
	}
	// fmt.Println(attributeCount, " ", occurenseCount)
	probabilities := make(map[string]float64) // probability for each decision attribute
	unAtribs := obj.getUniqueAttributes()
	uniqueAttr := unAtribs[:len(unAtribs)-1]        // Unique Attributes except last column
	uniqueDecisionAttr := unAtribs[len(unAtribs)-1] // Unique decision attribute
	for _, val := range uniqueDecisionAttr {
		probability := 1.0
		attsCountTmp := attributeCount[val]
		for k := 0; k < len(attsCountTmp); k++ {
			att := attsCountTmp[k] //Counted attributes for each decision attribute
			attProbability := 0.0
			if att != 0.0 {
				attProbability = float64(att) / float64(occurenseCount[val])
			} else {
				//smoothing
				attProbability = 1.0 / (float64(occurenseCount[val]) + float64(len(uniqueAttr[k])))
			}
			probability *= attProbability
		}
		probability *= float64(occurenseCount[val]) / float64(len(obj.Obj))
		probabilities[val] = probability
	}
	// fmt.Println(probabilities)

	// Finding the biggest probability value
	maxProb := 0.0
	maxProbName := ""
	for name, prob := range probabilities {
		if maxProb < prob {
			maxProb = prob
			maxProbName = name
		}
	}
	return maxProbName
}

func (obj *Objects) calculateAccuracy() float64 {
	tr := new(Objects)
	tr.readData("./data/train")
	correctAnsw := 0
	for ind, val := range obj.Obj {
		classifiedName := tr.classify(val.Attributes[:len(val.Attributes)-1])
		if classifiedName == obj.Obj[ind].Attributes[len(obj.Obj[ind].Attributes)-1] {
			correctAnsw++
		}
	}
	accurancy := (float64(correctAnsw) / float64(len(obj.Obj))) * 100.0
	return accurancy
}

func (obj *Objects) inputCheck() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the attributes: ")
	str, _ := reader.ReadString('\n')
	attrStr := strings.Fields(str) // Convert string to []string
	fmt.Println("\nClassified as:", obj.classify(attrStr))
}

func start() {
	ts := new(Objects)
	tr := new(Objects)
	test := "./data/test"
	ts.readData(test)
	tr.readData("./data/train")
	fmt.Println("\nAccuracy:", ts.calculateAccuracy(), "%")
	tr.inputCheck()
}
