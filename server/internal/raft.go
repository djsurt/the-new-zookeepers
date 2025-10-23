package election

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/djsurt/the-new-zookeepers/server/proto/raft"
	"google.golang.org/grpc"
)

type NodeState uint

const (
	FOLLOWER NodeState = iota
	CANDIDATE
	LEADER
)

type ElectionServer struct {
	raft.UnimplementedElectionServer
	Port       int
	state      NodeState
	grpcServer *grpc.Server
	listener   net.Conn
}

func NewElectionServer(port int) *ElectionServer {
	return &ElectionServer{
		Port:  port,
		state: FOLLOWER,
	}
}

// Handle a RequestVote call from a peer in the candidate state.
func (s *ElectionServer) RequestVote(
	ctx context.Context,
	req *raft.VoteRequest,
) (*raft.Vote, error) {
	log.Printf("Vote request received from %d", req.GetCandidateId())
	vote := &raft.Vote{Term: 1, VoteGranted: false}
	return vote, nil
}

// When in the leader state, make an an AppendEntries either to update a
// follower's log, or to send a heartbeat to the follower.
// When in the follower state, respond to AppendEntries requests and udpate
// election timeout.
func (s *ElectionServer) AppendEntries(
	ctx context.Context,
	req *raft.AppendEntriesRequest,
) (*raft.AppendEntriesResult, error) {
	log.Printf("Heartbeat received from %d", req.GetLeaderId())
	res := &raft.AppendEntriesResult{
		Term:    req.GetTerm(),
		Success: true,
	}
	return res, nil
}

// Serve the ElectionServer RPC interface. Returns an error if any of the
// setup steps fail, or if the grpcServer returns an error due to a
// listener.accept() failure.
func (s *ElectionServer) Serve() error {
	// Try to create the TCP socket.
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.Port))
	if err != nil {
		return fmt.Errorf("Error creating TCP socket: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	raft.RegisterElectionServer(s.grpcServer, s)

	// Begin serving ElectionServer RPCs
	err = s.grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}
