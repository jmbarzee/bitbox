package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	bbgrpc "github.com/jmbarzee/bitbox/grpc"
)

var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "query",
	Long:  "Query a process on the bitbox server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Require a single id as an argument")
		}

		uuid, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse uuid: %s", args[0])
		}

		job := jobQuery{
			id: uuid,
		}
		ctx := context.Background()
		bbClient := getClient(ctx)
		return job.execute(ctx, bbClient)
	},
}

type jobQuery struct {
	id uuid.UUID
}

// Execute querys a job on the remote bitBox
func (j jobQuery) execute(ctx context.Context, c bbgrpc.BitBoxClient) error {
	request := &bbgrpc.QueryRequest{
		ID: j.id[:],
	}

	queryClient, err := c.Query(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to stop process %s: %w", j.id, err)
	}

	for {
		reply, err := queryClient.Recv()
		if err == io.EOF {
			log.Println("<End of Stream>")
		}
		if err != nil {
			return fmt.Errorf("failed to fetch reply: %w", err)
		}
		switch output := reply.GetOutput().(type) {
		case *bbgrpc.QueryReply_Stdouterr:
			log.Println(output.Stdouterr)
		case *bbgrpc.QueryReply_ExitCode:
			log.Printf("Process %v exited with code %v", j.id, output.ExitCode)
			return nil
		}
	}
}
