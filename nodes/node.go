package nodes

import (
	"MP3/utils"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

func handleNode(wg *sync.WaitGroup, node utils.Node, config utils.Config) {
	defer wg.Done()
	valChan := make(chan utils.Message)
	go handleConnections(valChan, node.Server)
	processIncomingValues(valChan, node, config)
}

func handleConnections(valChan chan utils.Message, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go unicast_receive(valChan, conn)
	}
}

func unicast_receive(valChan chan utils.Message, conn net.Conn) {
	for {
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&message)
		utils.CheckError(err)
		valChan <- message
	}
}

func processIncomingValues(valChan chan utils.Message, node utils.Node, config utils.Config) {

	from := node.Id
	n := len(config.Nodes)
	f := config.F

	round := 1
	var value float64
	recievedMessages := 0
	sum := 0.
	canFail := true

	err := multicast(node.Conns,
		utils.Message{From: from, Round: round, Value: node.Input},
		true,
		config.MinDelay,
		config.MaxDelay)
	if err != nil {
		fmt.Println("Node "+node.Id+":", err.Error())
		unicast_send(node.Conns[0], utils.Message{From: node.Id, Fail: true}, 0, 1)
		return
	}

	for {
		message := <-valChan

		fmt.Printf("Node %s received message from %s for round %d\n", from, message.From, message.Round)

		if message.Output {
			println(value)
			break
		} else if message.Fail {
			canFail = false
		} else if message.Round < round {
			continue
		} else if message.Round > round {
			valChan <- message
		} else {
			recievedMessages++
			sum += message.Value

			if recievedMessages >= n-f {
				value = sum / float64(n-f)
				fmt.Printf("Node %s finished round %d ------------- \n", from, round)
				err = multicast(node.Conns,
					utils.Message{From: from, Round: round + 1, Value: value},
					canFail,
					config.MinDelay,
					config.MaxDelay)
				if err != nil {
					fmt.Println("Node "+node.Id+":", err.Error())
					unicast_send(node.Conns[0], utils.Message{From: node.Id, Fail: true}, 0, 1)
					return
				}
				round++
				recievedMessages = 0
				sum = 0.
			}
		}
	}
}

func multicast(conns []net.Conn, message utils.Message, canFail bool, min, max int) error {

	// Always send to master server
	encoder := gob.NewEncoder(conns[0])
	err := encoder.Encode(message)
	utils.CheckError(err)

	// Send to all other conns
	rand.Seed(time.Now().UnixNano())
	//intid, err := strconv.ParseInt(message.From, 10, 64)
	//rand.Seed(intid)
	for _, conn := range conns[1:] {
		if canFail {
			number := rand.Intn(100)
			if number < 3 {
				return errors.New("node has crashed")
			}
		}
		go unicast_send(conn, message, min, max)
	}

	return nil
}

func unicast_send(conn net.Conn, message utils.Message, min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(message)
	utils.CheckError(err)
}
