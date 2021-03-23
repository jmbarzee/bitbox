package bitbox

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmbarzee/bitbox/grpc"
	"github.com/jmbarzee/bitbox/proc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ensure Server implements BitBoxServer
var _ grpc.BitBoxServer = (*Server)(nil)

// Server starts, stops, and tracks arbitrary processes
type Server struct {
	// UnimplementedBitBoxServer is embedded to enable forwards compatability
	grpc.UnimplementedBitBoxServer
	c *Core
}

func NewServer() *Server {
	return &Server{
		c: NewCore(),
	}
}

// Start initiates a process.
func (s *Server) Start(ctx context.Context, request *grpc.StartRequest) (*grpc.StartReply, error) {
	cmd := request.GetCommand()
	args := request.GetArguments()
	log.Println("[Start] ", cmd, args)

	uuid, err := s.c.Start(cmd, args...)
	if err != nil {
		return nil, err
	}

	uuidBytes := uuid[:]

	return &grpc.StartReply{
		ID: uuidBytes,
	}, nil
}

// Stop halts a process.
func (s *Server) Stop(ctx context.Context, request *grpc.StopRequest) (*grpc.StopReply, error) {
	uuid, err := uuid.FromBytes(request.GetID())
	if err != nil {
		return nil, err
	}
	log.Println("[Stop] ", uuid.String())

	return &grpc.StopReply{}, s.c.Stop(uuid)
}

// Status returns the status of a process.
func (s *Server) Status(ctx context.Context, request *grpc.StatusRequest) (*grpc.StatusReply, error) {
	uuid, err := uuid.FromBytes(request.GetID())
	if err != nil {
		return nil, err
	}
	log.Println("[Status] ", uuid.String())

	status, err := s.c.Status(uuid)
	if err != nil {
		return nil, err
	}

	grpcStatus, err := convertToGRPCStatus(status)
	if err != nil {
		return nil, err
	}

	return &grpc.StatusReply{
		Status: grpcStatus,
	}, nil
}

func convertToGRPCStatus(status proc.ProcStatus) (grpc.StatusReply_StatusEnum, error) {
	// We could leverage implement this function on the proc.ProcStatus,
	// That would force the proc package to know import grpc, which is worth avoiding.
	switch status {
	case proc.Running:
		return grpc.StatusReply_Running, nil
	case proc.Exited:
		return grpc.StatusReply_Exited, nil
	case proc.Stopped:
		return grpc.StatusReply_Stopped, nil
	}
	return 0, fmt.Errorf("Unknown process status: %v", status)
}

// Query streams the output/result of a process.
func (s *Server) Query(request *grpc.QueryRequest, queryServer grpc.BitBox_QueryServer) error {
	uuid, err := uuid.FromBytes(request.GetID())
	if err != nil {
		return err
	}

	log.Println("[Query] ", uuid.String())
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
