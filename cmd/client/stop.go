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

var cmdStop = &cobra.Command{
	Use:   "stop",
	Short: "stop",
	Long:  "Stop a process on the bitbox server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic(errors.New("Require a single id as an argument"))
		}

		uuid, err := uuid.Parse(args[0])
		if err != nil {
			panic(err)
		}

		job := jobStop{
			id: uuid,
		}
		ctx := context.Background()
		bbClient := getClient(ctx)
		if err := job.execute(ctx, bbClient); err != nil {
			panic(err)
		}
	},
}

type jobStop struct {
	id uuid.UUID
}

// Execute stops a job on the remote BibBox
func (j jobStop) execute(ctx context.Context, c bbgrpc.BitBoxClient) error {
	request := &bbgrpc.StopRequest{
		ID: j.id[:],
	}

	_, err := c.Stop(ctx, request)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to stop process %s: %w", j.id, err))
	}
	log.Println("Successfully stopped process: ", j.id)
	return nil
}
