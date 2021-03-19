package bitbox

import (
	"context"

	"github.com/jmbarzee/bitbox/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ensure Server implements BitBoxServer
var _ grpc.BitBoxServer = (*Server)(nil)

// Server starts, stops, and tracks arbitrary processes
type Server struct {
	// UnimplementedBitBoxServer is embedded to enable forwards compatability
	grpc.UnimplementedBitBoxServer
}

func NewServer() *Server {
	return &Server{}
}

// Start initiates a process.
func (s *Server) Start(ctx context.Context, start *grpc.StartRequest) (*grpc.StartReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}

// Stop halts a process.
func (s *Server) Stop(context.Context, *grpc.StopRequest) (*grpc.StopReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}

// Status returns the status of a process.
func (s *Server) Status(context.Context, *grpc.StatusRequest) (*grpc.StatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}

// Query streams the output/result of a process.
func (s *Server) Query(*grpc.QueryRequest, grpc.BitBox_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
