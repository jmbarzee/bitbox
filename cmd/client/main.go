package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	bbgrpc "github.com/jmbarzee/bitbox/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var defaultPort = optionalEnvString("BIT_BOX_PORT", "8443")

var defaultAddress = optionalEnvString("BIT_BOX_ADDR", func() string {
	ip, err := getOutboundIP()
	if err != nil {
		panic(err)
	}
	return ip.String()
}())

var root = &cobra.Command{
	Use:   "bitboxc",
	Short: "A CLI tool for remote BitBox operations",
	Long:  "Execute remote linux processes on a BitBox server",
}

func main() {
	root.AddCommand(cmdStart)
	root.AddCommand(cmdStop)
	root.AddCommand(cmdStatus)
	root.AddCommand(cmdQuery)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}

func optionalEnvString(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Missing optional environment variable %s, using default %s", key, defaultValue)
	return defaultValue
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func getClient(ctx context.Context) bbgrpc.BitBoxClient {
	address := fmt.Sprintf("%s:%s", defaultAddress, defaultPort)

	conn, err := grpc.DialContext(
		context.TODO(),
		address,
		//TODO: replace with mTLS
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		panic(fmt.Errorf("Failed to dial connection during reconnect: %w", err))
	}

	return bbgrpc.NewBitBoxClient(conn)
}
