package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	bbgrpc "github.com/jmbarzee/bitbox/grpc"
)

var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "query",
	Long:  "Query a process on the bitbox server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic(errors.New("Require a single id as an argument"))
		}

		job := jobQuery{
			id: args[0],
		}
		ctx := context.Background()
		bbClient := getClient(ctx)
		if err := job.execute(ctx, bbClient); err != nil {
			panic(err)
		}
	},
}

type jobQuery struct {
	id string
}

// Execute querys a job on the remote BibBox
func (j jobQuery) execute(ctx context.Context, c bbgrpc.BitBoxClient) error {
	request := &bbgrpc.QueryRequest{
		ID: []byte(j.id),
	}

	queryClient, err := c.Query(ctx, request)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to stop process %s: %w", j.id, err))
	}

Loop:
	for {
		reply, err := queryClient.Recv()
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch reply: %w", err))
			log.Println(fmt.Errorf("failed to fetch reply: %w", err))
			break
		}
		switch output := reply.GetOutput().(type) {
		case *bbgrpc.QueryReply_Stdout:
			log.Println(output.Stdout)
		case *bbgrpc.QueryReply_Stderr:
			log.Printf("Error: %v", output.Stderr)
		case *bbgrpc.QueryReply_ExitCode:
			log.Printf("Process %v exited with code %v", j.id, output.ExitCode)
			break Loop
		}
	}

	return nil
}
