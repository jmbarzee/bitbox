package bitbox

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative bitbox.proto

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmbarzee/bitbox/grpc"
	"github.com/jmbarzee/bitbox/proc"
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

	return &grpc.StartReply{
		ID: uuid[:],
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
	// TODO: pass context from queryServer to Query
	stream, err := s.c.Query(queryServer.Context(), uuid)
	if err != nil {
		return err
	}

	for output := range stream {
		// You should ask yourself: why do we initialize reply the same way twice?
		// The answer: because QueryReply.Output is of type grpc.isQueryReply_Output
		// which is not public.
		var reply *grpc.QueryReply
		switch output := output.(type) {
		case *proc.ProcOutput_Stdouterr:
			reply = &grpc.QueryReply{
				Output: &grpc.QueryReply_Stdouterr{
					Stdouterr: output.Output,
				},
			}
		case *proc.ProcOutput_ExitCode:
			reply = &grpc.QueryReply{
				Output: &grpc.QueryReply_ExitCode{
					ExitCode: output.ExitCode,
				},
			}
		}
		err := queryServer.Send(reply)
		if err != nil {
			return err // TODO: is this really how we should handle this?
		}
	}
	return nil
}
