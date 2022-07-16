package raft

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/FlorinBalint/flosensus/raft/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Role int

const (
	Stopped Role = iota
	Initialized
	Follower
	Candidate
	Leader
	Snapshotting // TODO(#6): Implement log compaction
)

const (
	notVoted = -1
)

// solidState is the persisted state for each server.
// All updates are persisted before responding to a request
// (a vote or a new log entry)
type solidState struct {
	// last seen Term
	currentTerm int64
	// candidate that received the vote in current term (-1 if none)
	votedFor int32
	// log entries, each entry contains the command for state machine,
	// and the term received from the leader
	logs []pb.LogEntry
}

func (s solidState) Persist() {
	//TODO(#3): Persist the necessary state
}

// leaderState is the leaders knowledge / view about the system.
type leaderState struct {
	// the log index to send to each server
	// initial value: the last leader log index +1
	nextIndex []int64
	// the highest log entry known to be replicated by each server
	// initial value: 0
	matchindex []int64
}

// volatileState is the state that is not persisted during runtime
type volatileState struct {
	// index of highest commited known log entry (initially 0)
	commitIndex int64
	// index of the highest applied (not necessarily commited)
	// log entry (initially 0)
	lastApplied int64
	// Additional information for leader, reinitialized after election
	leader leaderState
}

// state is the current state of the server, including both
// volatile and non-volatile info.
type state struct {
	solid    solidState
	volatile volatileState
}

type peer struct {
	id     int
	client pb.RaftClient
}

// server is used to implement flosensus.Server.
type Server struct {
	pb.UnimplementedRaftServer

	grpcServer *grpc.Server
	currRole   Role
	currState  state

	id         int
	port       int
	minTimeout time.Duration
	maxTimeout time.Duration
	peers      map[int]peer

	mutex sync.RWMutex
}

func (rs *Server) AppendEntries(ctx context.Context, in *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	// TODO(#1): Implement log replication
	return nil, nil
}

func (rs *Server) RequestVote(ctx context.Context, in *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	// TODO(#2): Implement leader election
	return nil, nil
}

func (rs *Server) role() Role {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.currRole
}

func (rs *Server) followerLoop() {
	// TODO(#1): Implement log replication
}

func (rs *Server) candidateLoop() {
	// TODO(#2): Implement leader election
}

func (rs *Server) leaderLoop() {
	// TODO(#1 & #2): Implement leader election & log replication
}

func (rs *Server) snapshotLoop() {
	// TODO(#6): Implement log compaction
}

func (rs *Server) activeLoop() {
	role := rs.role()

	for role != Stopped {
		switch role {
		case Follower:
			rs.followerLoop()
		case Candidate:
			rs.candidateLoop()
		case Leader:
			rs.leaderLoop()
		case Snapshotting:
			rs.snapshotLoop()
		}
		role = rs.role()
	}
}

func initPeers(peerCfgs []*pb.PeerConfig) (map[int]peer, error) {
	res := make(map[int]peer)
	for _, peerCfg := range peerCfgs {
		addr := fmt.Sprintf("%v:%v", peerCfg.GetAddress(), peerCfg.GetPort())
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("could not create peer connection: %v", err)
		}
		peer := peer{
			id:     int(peerCfg.GetId()),
			client: pb.NewRaftClient(conn),
		}
		res[peer.id] = peer
	}
	return res, nil
}

func initState(cfg *pb.Config) state {
	solid := solidState{
		votedFor: notVoted,
	}
	solid.Persist()
	return state{
		solid: solid,
	}
}

func NewServer(cfg *pb.Config) (*Server, error) {
	peers, err := initPeers(cfg.GetPeers())
	if err != nil {
		return nil, err
	}
	server := &Server{
		id:         int(cfg.GetSelf().GetId()),
		port:       int(cfg.GetSelf().GetPort()),
		minTimeout: cfg.GetMinElectionTimeout().AsDuration(),
		maxTimeout: cfg.GetMaxElectionTimeout().AsDuration(),
		peers:      peers,
		currState:  initState(cfg),
	}
	return server, nil
}

func (rs *Server) ListenAndServe(_ context.Context) error {
	rs.currRole = Follower
	rs.grpcServer = grpc.NewServer()
	pb.RegisterRaftServer(rs.grpcServer, rs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", rs.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Raft Server will listen on %v", lis.Addr())

	go func() {
		rs.activeLoop()
	}()

	return rs.grpcServer.Serve(lis)
}
