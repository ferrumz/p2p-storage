package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

type node struct {
	Addr     string
	Port     string
	Status   string
	LastCall time.Time
}

type storageValue struct {
	Value string
	Time  time.Time
}

var nodes = map[string]node{}
var storage = map[string]storageValue{}

var mtx sync.Mutex

func main() {

	ip := flag.String("ip", "127.0.0.1", "ip address")
	port := flag.String("port", "2222", "port")
	masterNodeIp := flag.String("nodeIp", "", "master node ip")
	masterNodePort := flag.String("nodePort", "", "master node port")
	flag.Parse()
	nodes[*ip+":"+*port] = node{
		Addr:   *ip,
		Port:   *port,
		Status: "online",
	}

	var err error

	log.Print("server starting...")

	// sync with masternode if exists
	if *masterNodeIp != "" && *masterNodePort != "" {
		nodes[*masterNodeIp+":"+*masterNodePort] = node{
			Addr:   *masterNodeIp,
			Port:   *masterNodePort,
			Status: "online",
		}
		syncWithNode(*masterNodeIp + ":" + *masterNodePort)
	}

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
			}

			go func() {
				for {
					message, _ := bufio.NewReader(conn).ReadString('\n')
					message = strings.Trim(message, " \n\r")
					fmt.Print("Message: ", string(message), "\n")
					var response []byte
					switch {
					case message == "GET NODES":
						response, err = json.Marshal(nodes)
						if err != nil {
							log.Fatal(err)
						}
					case func() bool {
						matched, ok := regexp.MatchString(`GET KEY \w+`, message)
						return ok == nil && matched == true
					}():
						var exp = regexp.MustCompile(`GET KEY (\w+)`)
						key := exp.FindStringSubmatch(message)
						mtx.Lock()
						response = []byte(storage[key[1]].Value)
						mtx.Unlock()
						if err != nil {
							log.Fatal(err)
						}
					case func() bool {
						matched, ok := regexp.MatchString(`SET KEY \w+ \w+`, message)
						return ok == nil && matched == true
					}():
						var exp = regexp.MustCompile(`SET KEY (\w+) (\w+)`)
						key := exp.FindStringSubmatch(message)
						mtx.Lock()
						storage[key[1]] = storageValue{
							Value: key[2],
							Time:  time.Now(),
						}
						mtx.Unlock()
						response = []byte("Done")
						if err != nil {
							log.Fatal(err)
						}
					case message == "GET KEYS":
						mtx.Lock()
						response, err = json.Marshal(storage)
						mtx.Unlock()
						if err != nil {
							log.Fatal(err)
						}
					default:
						response = []byte("Unknown command")
					}
					_, err = conn.Write(append(response, []byte("\n")...))
					if err != nil {
						log.Print(err)
						return
					}
				}
			}()
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		for _, n := range nodes {
			if n.LastCall.Before(time.Now().Add(-5 * time.Second)) {
				syncWithNode(n.Addr + ":" + n.Port)
				nodes[n.Addr+":"+n.Port] = node{n.Addr, n.Port, n.Status, time.Now()}
			}
		}
	}

}

func getFromNode(node string, command string) (string, error) {
	var err error
	conn, err := net.Dial("tcp", node)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(conn, command+"\n")
	if err != nil {
		return "", err
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	if conn != nil {
		conn.Close()
	}
	return message, nil
}

func syncWithNode(nodeAddr string) {
	var remoteNodes map[string]node
	var message string
	var err error

	message, err = getFromNode(nodeAddr, "GET NODES")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte(message), &remoteNodes)
	if err != nil {
		log.Fatal(err)
	}
	syncNodes(remoteNodes)

	var remoteStorage map[string]storageValue
	message, err = getFromNode(nodeAddr, "GET KEYS")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte(message), &remoteStorage)
	if err != nil {
		log.Fatal(err)
	}
	syncStorage(remoteStorage)
}

func syncNodes(nodes map[string]node) {
	for _, remoteNode := range nodes {
		if nodes[remoteNode.Addr+":"+remoteNode.Port] == (node{}) {
			nodes[remoteNode.Addr+":"+remoteNode.Port] = remoteNode
		}
	}
}

func syncStorage(storageValues map[string]storageValue) {
	for _, value := range storageValues {
		if storage[value.Value] != (storageValue{}) {
			if value.Time.After(storage[value.Value].Time) {
				storage[value.Value] = value
			}
		} else {
			storage[value.Value] = value
		}
	}
}
