package cmd

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	otelProfileService "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"
	"google.golang.org/grpc"
)

type Server struct {
	otelProfileService.UnimplementedProfilesServiceServer
}

func (s Server) Export(ctx context.Context, req *otelProfileService.ExportProfilesServiceRequest) (*otelProfileService.ExportProfilesServiceResponse, error) {
	t := time.Now()
	fmt.Printf("\n------------------------------------------\n")
	fmt.Printf("Received at %d:%d:%d\n\n", t.Hour(), t.Minute(), t.Second())
	PrettyPrint(req)
	response := &otelProfileService.ExportProfilesServiceResponse{}
	return response, nil
}

var serverCmd = &cobra.Command{
	Use:   "server [port]",
	Args:  cobra.ExactArgs(1),
	Short: "Start a gRPC server waiting for OTEL profile data",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			return err
		}

		server := grpc.NewServer()
		otelProfileService.RegisterProfilesServiceServer(server, Server{})
		fmt.Printf("Listening for OTEL profiles via gRPC on port %d...\n", port)

		return server.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
