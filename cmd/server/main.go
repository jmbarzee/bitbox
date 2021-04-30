package main

import (
	// temp import for build
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jmbarzee/bitbox"
	bbgrpc "github.com/jmbarzee/bitbox/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var defaultPort = optionalEnvString("BIT_BOX_PORT", "8443")
var defaultAddress = optionalEnvString("BIT_BOX_ADDR", "")

func main() {

	address := fmt.Sprintf("%s:%s", defaultAddress, defaultPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Errorf("failed to listen on %s: %w", address, err))
	}

	bitBoxServer := bitbox.NewServer()

	server := grpc.NewServer()

	bbgrpc.RegisterBitBoxServer(server, bitBoxServer)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	err = server.Serve(lis)
}

func optionalEnvString(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Missing optional environment variable %s, using default %s", key, defaultValue)
	return defaultValue
}
