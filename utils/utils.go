package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	MinDelay int
	MaxDelay int
	F        int
	Nodes    []Node
	MServer  Server
}

type Node struct {
	Id     string
	Input  float64
	Ip     string
	Port   string
	Conns  []net.Conn
	Server net.Listener
}

type Server struct {
	Ip    string
	Port  string
	Conns []net.Conn
}

type Message struct {
	From  string
	Round int
	Value float64
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

	err = file.Close()
	CheckError(err)

	// Get min delay, max delay, and f
	line := strings.Split(txtlines[0], " ")
	min, _ := strconv.Atoi(line[0])
	max, _ := strconv.Atoi(line[1])
	f, _ := strconv.Atoi(line[2])

	// Get master server info
	line = strings.Split(txtlines[1], " ")
	server := Server{Ip: line[0], Port: line[1], Conns: []net.Conn{}}

	// Get list of nodes. Loop through config file lines, skipping line 1 since it contains delay params
	var nodeList []Node
	for _, line := range txtlines[2:] {
		// For each line, create node struct and add it to list of nodes
		list := strings.Split(line, " ")
		input, err := strconv.ParseFloat(list[1], 64)
		CheckError(err)
		node := Node{Id: list[0], Input: input, Ip: list[2], Port: list[3], Conns: []net.Conn{}}
		nodeList = append(nodeList, node)
	}

	return Config{MinDelay: min, MaxDelay: max, F: f, Nodes: nodeList, MServer: server}
}
