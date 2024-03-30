package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/fokus/proto"
)

type Fokusable interface {
	getFokus(ctx context.Context) (*pb.Focus, error)
	getName() string
	getType() pb.Focus_FocusType
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetFokus(ctx context.Context, req *pb.GetFokusRequest) (*pb.GetFokusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Need to get to this")
}
