package main

import (
	"bufio"
	"fmt"
	"github.com/hashicorp/raft"
	"os"
	"strconv"
	"strings"
	"time"
)

type RaftCluster struct {
	nodes []*raft.Raft
}

func NewRaftCluster(size int) (*RaftCluster, error) {
	nodes, err := raftInit(size)
	if err != nil {
		return nil, err
	}
	return &RaftCluster{nodes: nodes}, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var cluster *RaftCluster

	for {
		fmt.Print("Enter Quorum Members: ")
		if scanner.Scan() {
			input := scanner.Text()
			num, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input, please enter a number.")
				continue
			}
			cluster, err = NewRaftCluster(num)
			if err != nil {
				fmt.Println("RaftInit Err.")
				break
			}
			fmt.Println("Quorum initialized with", num, "members.")
			break
		}
	}

	scanner = bufio.NewScanner(os.Stdin)
	time.Sleep(2 * time.Second)
	for {
		fmt.Print("Enter command \n 1.kill node \n 2.list \n 3.exit \n: ")
		if !scanner.Scan() {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		handleCommand(line, cluster)
	}

	select {}
}

func handleCommand(input string, cluster *RaftCluster) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "kill":
		handleKillCommand(parts, cluster)
		time.Sleep(4 * time.Second)
	case "list":
		SelectAllNode(cluster)
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Unknown command")
	}
}

func handleKillCommand(parts []string, cluster *RaftCluster) {
	if len(parts) != 2 {
		fmt.Println("Usage: kill nodeID")
		return
	}
	nodeID := parts[1]

	for i, v := range cluster.nodes {

		var name = fmt.Sprintf("node%d", i+1)
		if name == nodeID {
			killNode(v)
			break
		}
	}
}

func killNode(raftNode *raft.Raft) {
	if raftNode == nil {
		fmt.Println("Cannot kill node: node is nil")
		return
	}
	raftNode.Shutdown()
	fmt.Printf("Node has been shutdown.\n")
}

func SelectAllNode(cluster *RaftCluster) {
	for i, v := range cluster.nodes {
		fmt.Printf("node%v: %v .\n", i+1, v)
	}
}
