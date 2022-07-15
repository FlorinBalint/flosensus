package raft

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/FlorinBalint/flosensus/raft/proto"
)

type State int

const (
	Stopped State = iota
	Initialized
	Follower
	Candidate
	Leader
	Snapshotting // TODO(#6): Implement log compaction
)

// server is used to implement flosensus.RaftServer.
type RaftServer struct {
	pb.UnimplementedRaftServer

	currState State
	cfg       *pb.Config
	mutex     sync.RWMutex
}

func (rs *RaftServer) AppendEntries(ctx context.Context, in *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	// TODO(#1): Implement log replication
	return nil, nil
}

func (rs *RaftServer) RequestVote(ctx context.Context, in *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	// TODO(#2): Implement leader election
	return nil, nil
}

func (rs *RaftServer) state() State {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.currState
}

func (rs *RaftServer) followerLoop() {
	// TODO(#1): Implement log replication
}

func (rs *RaftServer) candidateLoop() {
	// TODO(#2): Implement leader election
}

func (rs *RaftServer) leaderLoop() {
	// TODO(#1 & #2): Implement leader election & log replication
}

func (rs *RaftServer) snapshotLoop() {
	// TODO(#6): Implement log compaction
}

func (rs *RaftServer) activeLoop() {
	state := rs.state()

	for state != Stopped {
		switch state {
		case Follower:
			rs.followerLoop()
		case Candidate:
			rs.candidateLoop()
		case Leader:
			rs.leaderLoop()
		case Snapshotting:
			rs.snapshotLoop()
		}
		state = rs.state()
	}
}

func NewServer(cfg *pb.Config) *RaftServer {
	server := &RaftServer{
		cfg: cfg,
	}
	return server
}

func (rs *RaftServer) Start() error {
	rs.currState = Follower

	go func() {
		rs.activeLoop()
	}()
	return fmt.Errorf("Not implemented!")
}
