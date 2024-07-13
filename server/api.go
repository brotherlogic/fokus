package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ghbclient "github.com/brotherlogic/githubridge/client"
	githubridgeclient "github.com/brotherlogic/githubridge/client"

	pb "github.com/brotherlogic/fokus/proto"
)

type Fokusable interface {
	getFokus(ctx context.Context, client githubridgeclient.GithubridgeClient, now time.Time) (*pb.Focus, error)
	getName() string
	getType() pb.Focus_FocusType
}

type Server struct {
	modules []Fokusable
	client  githubridgeclient.GithubridgeClient
}

func NewServer() *Server {
	client, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to reach GHB")
	}
	return &Server{
		modules: []Fokusable{&RecordAdd{}, &Highlight{}, &HomeWeek{}, &Home{}, &Overdue{}},
		client:  client,
	}
}

func (s *Server) GetFokus(ctx context.Context, req *pb.GetFokusRequest) (*pb.GetFokusResponse, error) {
	// We can't rely on America/Los_Angeles being present it seems; ignore Daylight savbings
	location := time.FixedZone("UTC-8", -7*60*60)
	t := time.Now().In(location)

	for _, m := range s.modules {
		focus, err := m.getFokus(ctx, s.client, t)
		log.Printf("%v -> %v", m.getName(), err)
		if err == nil && focus != nil {
			return &pb.GetFokusResponse{Focus: focus}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "Could not find focus task")
}
