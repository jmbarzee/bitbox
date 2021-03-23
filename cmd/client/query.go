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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic(errors.New("Require a single id as an argument"))
		}

		uuid, err := uuid.Parse(args[0])
		if err != nil {
			panic(err)
		}

		job := jobQuery{
			id: uuid,
		}
		ctx := context.Background()
		bbClient := getClient(ctx)
		if err := job.execute(ctx, bbClient); err != nil {
			panic(err)
		}
	},
}

type jobQuery struct {
	id uuid.UUID
}

// Execute querys a job on the remote BibBox
func (j jobQuery) execute(ctx context.Context, c bbgrpc.BitBoxClient) error {
	request := &bbgrpc.QueryRequest{
		ID: j.id[:],
	}

	queryClient, err := c.Query(ctx, request)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to stop process %s: %w", j.id, err))
	}

	for {
		reply, err := queryClient.Recv()
		if err == io.EOF {
			log.Println("<End of Stream>")
			break
		}
		if err != nil {
			log.Fatal(fmt.Errorf("failed to fetch reply: %w", err))
		}
		log.Print(reply.GetOutput())
	}

	return nil
}
