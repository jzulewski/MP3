package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	MinDelay int
	MaxDelay int
	Nodes    []Node
}

type Node struct {
	Id   string
	Ip   string
	Port string
}

// Consolidated repeated error checks into a single function
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// TODO: Fix potential indexing error
func CreateConfigStruct() Config {
	file, err := os.Open("config.txt")
	CheckError(err)

	// Create scanner object and textlines array
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	// Loop through file lines, appending to textlines
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()

	// Get min and max delay
	line := strings.Split(txtlines[0], " ")
	min, _ := strconv.Atoi(line[0])
	max, _ := strconv.Atoi(line[1])

	// Get list of nodes. Loop through config file lines, skipping line 1 since it contains delay params
	var nodeList []Node
	for _, line := range txtlines[1:] {
		// For each line, create node struct and add it to list of nodes
		list := strings.Split(line, " ")
		node := Node{Id: list[0], Ip: list[1], Port: list[2]}
		nodeList = append(nodeList, node)
	}

	return Config{MinDelay: min, MaxDelay: max, Nodes: nodeList}
}
