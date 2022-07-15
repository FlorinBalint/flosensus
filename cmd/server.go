package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/FlorinBalint/flosensus/raft"
	pb "github.com/FlorinBalint/flosensus/raft/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	raftServer := raft.NewServer(nil)
	_ = raftServer.Start() // ignore error for now
	s := grpc.NewServer()
	pb.RegisterRaftServer(s, raftServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
