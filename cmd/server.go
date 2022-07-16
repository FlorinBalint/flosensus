package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FlorinBalint/flosensus/raft"
	pb "github.com/FlorinBalint/flosensus/raft/proto"
	"google.golang.org/protobuf/encoding/prototext"
)

var (
	config     = flag.String("config", "", "Config proto to use for the load balancer")
	configFile = flag.String("config_file", "", "Config file to use for the load balancer")
	port       = flag.Int("port", 50051, "The server port")
)

func parseConfig(cfg []byte) (*pb.Config, error) {
	res := &pb.Config{}
	err := prototext.Unmarshal(cfg, res)
	overridePortIfNeeded(res)
	return res, err
}

func parseConfigFromFile(cfgPath string) (*pb.Config, error) {
	content, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	return parseConfig(content)
}

func overridePortIfNeeded(cfg *pb.Config) {
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			*cfg.GetSelf().Port = int32(*port)
		}
	})
}

func main() {
	flag.Parse()
	if (len(*config) == 0) == (len(*configFile) == 0) {
		log.Fatalf("Exactly one of --config or --config_file must pe specified !")
	}

	var err error
	var raftCfg *pb.Config
	if len(*config) != 0 {
		raftCfg, err = parseConfig([]byte(*config))
	} else if len(*configFile) != 0 {
		raftCfg, err = parseConfigFromFile(*configFile)
	}

	if err != nil {
		log.Fatalf("Error while parsing the configs: %v\n", err)
	}

	raftServer, err := raft.NewServer(raftCfg)
	if err != nil {
		log.Fatalf("Could not create server: %v", err)
	}
	err = raftServer.ListenAndServe(context.Background())
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Raft server was stopped")
	} else if err != nil {
		fmt.Printf("Unexpected raft server error: %v\n The server shut down ...", err)
	}
}
