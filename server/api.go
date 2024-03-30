package server

import (
	"context"

	pb "github.com/brotherlogic/fokus/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	return &Server{
		modules: []Fokusable{&Overdue{}},
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
