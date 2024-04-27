package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"net"
	"os"
	"path/filepath"
	"time"
)

func setupLogger() hclog.Logger {
	logFile, err := os.OpenFile("log/raft.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "raft",
		Level:  hclog.LevelFromString("DEBUG"), // or another appropriate level
		Output: logFile,
	})
	return logger
}

func setupRaftNode(id, dataDir, raftAddress string) (*raft.Raft, *raft.NetworkTransport) {
	logger := setupLogger()
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(id)
	config.Logger = logger

	addr, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		logger.Error("Unable to resolve raft address", "address", raftAddress, "error", err)
		os.Exit(1)
	}

	transport, err := raft.NewTCPTransport(raftAddress, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		logger.Error("Unable to create raft transport", "error", err)
		os.Exit(1)
	}

	logsPath := filepath.Join(dataDir, "raft", "logs")
	os.MkdirAll(logsPath, 0755)
	logStore, err := boltdb.NewBoltStore(filepath.Join(logsPath, "raft.db"))
	if err != nil {
		logger.Error("Unable to create log store", "path", filepath.Join(logsPath, "raft.db"), "error", err)
		os.Exit(1)
	}

	stableStore, err := boltdb.NewBoltStore(filepath.Join(dataDir, "raft", "stable.db"))
	if err != nil {
		logger.Error("Unable to create stable store", "path", filepath.Join(dataDir, "raft", "stable.db"), "error", err)
		os.Exit(1)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(dataDir, 1, nil)
	if err != nil {
		logger.Error("Unable to create snapshot store", "path", dataDir, "error", err)
		os.Exit(1)
	}

	r, err := raft.NewRaft(config, nil, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		logger.Error("Failed to create raft", "error", err)
		os.Exit(1)
	}

	return r, transport
}

func raftInit(nodeNum int) ([]*raft.Raft, error) {
	logger := setupLogger()
	config := raft.DefaultConfig()
	config.Logger = logger
	var rafts = make([]*raft.Raft, nodeNum)
	var servers []raft.Server
	for i := 1; i <= nodeNum; i++ {
		dataDirNode := fmt.Sprintf("./raftData/node%d", i)
		os.RemoveAll(dataDirNode)
		os.MkdirAll(dataDirNode, 0755)

		node, transport := setupRaftNode(fmt.Sprintf("node%d", i), dataDirNode, fmt.Sprintf("127.0.0.1:500%d", i))

		server := raft.Server{
			ID:      raft.ServerID(fmt.Sprintf("node%d", i)),
			Address: transport.LocalAddr(),
		}
		servers = append(servers, server)
		rafts[i-1] = node
		go monitorLeadership(node)
	}

	configuration := raft.Configuration{Servers: servers}
	if len(rafts) > 0 {
		future := rafts[0].BootstrapCluster(configuration)
		if err := future.Error(); err != nil {
			logger.Error("Failed to bootstrap cluster", "error", err)
			return nil, err
		}
	}

	return rafts, nil
}

func monitorLeadership(node *raft.Raft) {
	for {
		select {
		case isLeader := <-node.LeaderCh():
			if isLeader {
				fmt.Printf("%s is now the leader\n", node)
			}
		}
	}
}
