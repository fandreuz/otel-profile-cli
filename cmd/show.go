package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	otelProfile "go.opentelemetry.io/proto/otlp/profiles/v1development"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

var indent string

func PrettyPrint(pb proto.Message) {
	out, err := prototext.MarshalOptions{
		Multiline: true,
		Indent:    indent,
	}.Marshal(pb)
	if err != nil {
		fmt.Println("Failed to pretty print:", err)
		return
	}
	fmt.Println(string(out))
}

var showCmd = &cobra.Command{
	Use:   "show [file]",
	Args:  cobra.ExactArgs(1),
	Short: "Pretty-print the content of a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		var profilesData otelProfile.ProfilesData
		if err := proto.Unmarshal(file, &profilesData); err != nil {
			return err
		}
		PrettyPrint(&profilesData)

		return nil
	},
}

func init() {
	showCmd.Flags().StringVar(&indent, "indent", "   ", "Indentation")
	rootCmd.AddCommand(showCmd)
}
