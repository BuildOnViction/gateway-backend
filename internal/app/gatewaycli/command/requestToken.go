package command

import (
	"context"
	"fmt"
	"time"

	gateway "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type createOptions struct {
	address string
	client  gateway.AuthServiceClient
}

// NewRequestTokenCommand creates a new cobra.Command for adding a new item to the list.
func NewRequestTokenCommand(c Context) *cobra.Command {
	options := createOptions{}

	cmd := &cobra.Command{
		Use:     "request-token",
		Aliases: []string{"rt"},
		Short:   "Request Login token",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.address = args[0]
			options.client = c.GetAuthServiceClient()

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runRequestToken(options)
		},
	}

	return cmd
}

func runRequestToken(options createOptions) error {
	req := &gateway.RequestTokenRequest{
		Address: options.address,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := options.client.RequestToken(ctx, req)
	if err != nil {
		st := status.Convert(err)
		for _, detail := range st.Details() {
			// nolint: gocritic
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				fmt.Println("Oops! Your request was rejected by the server.")
				for _, violation := range t.GetFieldViolations() {
					fmt.Printf("The %q field was wrong:\n", violation.GetField())
					fmt.Printf("\t%s\n", violation.GetDescription())
				}
			}
		}

		return err
	}

	fmt.Printf("Issued Token for logging in .", options.address, resp.Token)

	return nil
}
