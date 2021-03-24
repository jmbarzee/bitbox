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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Require a single id as an argument")
		}

		uuid, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("failed to parse uuid: %s", args[0])
		}

		job := jobStatus{
			id: uuid,
		}
		ctx := context.Background()
		bbClient := getClient(ctx)
		return job.execute(ctx, bbClient)
	},
}

type jobStatus struct {
	id uuid.UUID
}

// Execute returns the status of a job on the remote bitBox
func (j jobStatus) execute(ctx context.Context, c bbgrpc.BitBoxClient) error {
	request := &bbgrpc.StatusRequest{
		ID: j.id[:],
	}

	reply, err := c.Status(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to query process %s: %w", j.id, err)
	}
	log.Println("Successfully queried status of process: ", j.id, ", ", reply.Status.String())

	return nil
}
