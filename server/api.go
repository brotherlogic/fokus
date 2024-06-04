package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ghbclient "github.com/brotherlogic/githubridge/client"

	pb "github.com/brotherlogic/fokus/proto"
)

type Fokusable interface {
	getFokus(ctx context.Context) (*pb.Focus, error)
	getName() string
	getType() pb.Focus_FocusType
}

type Server struct {
	modules []Fokusable
}

func NewServer() *Server {
	client, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to reach GHB")
	}
	return &Server{
		modules: []Fokusable{&HomeWeek{client: client}, &Home{client: client}, &Overdue{client: client}},
	}
}

func (s *Server) GetFokus(ctx context.Context, req *pb.GetFokusRequest) (*pb.GetFokusResponse, error) {
	for _, m := range s.modules {
		focus, err := m.getFokus(ctx)
		if err == nil && focus != nil {
			return &pb.GetFokusResponse{Focus: focus}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "Could not find focus task")
}
