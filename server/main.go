package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/djsurt/the-new-zookeepers/server/proto/raft"
	"google.golang.org/grpc"
)

type ElectionServer struct {
	raft.UnimplementedElectionServer
	clientID int
}

func (s *ElectionServer) RequestVote(ctx context.Context, req *raft.VoteRequest) (*raft.Vote, error) {
	log.Printf("Vote request received from %d", req.GetCandidateId())
	vote := &raft.Vote{Term: 1, VoteGranted: false}
	return vote, nil
}

func main() {
	var port int
	flag.IntVar(&port, "port", -1, "The port number to use for this node")
	flag.Parse()

	// TODO: Check that the supplied port is in the cluster config.
	if port < 0 {
		fmt.Printf("Please provide a port number to bind to that is in the cluster file.\n")
		os.Exit(1)
	}

	_, err := parseClusterConfig(port, "cluster.conf")
	if err != nil {
		fmt.Printf("Error reading cluster config: %v\n", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Printf("Failed to create tcp listener: %v\n", err)
		os.Exit(1)
	}

	electionServer := &ElectionServer{
		clientID: port,
	}

	grpcServer := grpc.NewServer()
	raft.RegisterElectionServer(grpcServer, electionServer)

	serverErrChan := make(chan error)
	go func(errChan chan<- error) {
		err := grpcServer.Serve(listener)
		if err != nil {
			errChan <- err
		}
	}(serverErrChan)

	fmt.Printf("Election server listening on port %d...\n", port)
	// Wait for an error
	err = <-serverErrChan
	fmt.Printf("Error from election server: %v\n", err)
	os.Exit(1)
}

func parseClusterConfig(myPort int, configPath string) (peers []int, err error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	for line := range strings.SplitSeq(string(bytes), "\n") {
		if line == "" {
			continue
		}
		peer, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Ignoring entry %s: could not parse to integer\n", line)
			continue
		}
		if peer != myPort {
			peers = append(peers, peer)
		}
	}
	return peers, nil
}
