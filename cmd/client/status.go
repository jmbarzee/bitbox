package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	bbgrpc "github.com/jmbarzee/bitbox/grpc"
)

var cmdStatus = &cobra.Command{
	Use:   "status",
	Short: "status",
	Long:  "Stop a process on the bitbox server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic(errors.New("Require a single id as an argument"))
		}

		uuid, err := uuid.Parse(args[0])
		if err != nil {
			panic(err)
		}

		job := jobStatus{
			id: uuid,
		}
		ctx := context.Background()
		bbClient := getClient(ctx)
		if err := job.execute(ctx, bbClient); err != nil {
			panic(err)
		}
	},
}

type jobStatus struct {
	id uuid.UUID
}

// Execute returns the status of a job on the remote BibBox
func (j jobStatus) execute(ctx context.Context, c bbgrpc.BitBoxClient) error {
	request := &bbgrpc.StatusRequest{
		ID: j.id[:],
	}

	reply, err := c.Status(ctx, request)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to stop process %s: %w", j.id, err))
	}
	log.Println("Successfully queried status of process: ", j.id, ", ", reply.Status.String())

	return nil
}
