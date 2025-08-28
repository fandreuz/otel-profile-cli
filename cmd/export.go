/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	otelProfileService "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [file] [endpoint]",
	Args:  cobra.ExactArgs(2),
	Short: "Export the given profile with a gRPC request to the given endpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		var request otelProfileService.ExportProfilesServiceRequest
		if err := proto.Unmarshal(file, &request); err != nil {
			return err
		}

		conn, err := grpc.NewClient(args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()
		client := otelProfileService.NewProfilesServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		if response, err := client.Export(ctx, &request); err != nil {
			return err
		} else {
			PrettyPrint(response)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
