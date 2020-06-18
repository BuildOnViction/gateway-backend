package command

import (
	"context"
	"fmt"
	"time"

	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type listingProjectOptions struct {
	accessToken string
	authClient  gateway.AuthServiceClient
	projectCl   gateway.ProjectServiceClient
}

// NewListprojectCommand lists a new cobra.Command for adding a new item to the list.
func NewListProjectCommand(c Context) *cobra.Command {
	options := listingProjectOptions{}

	cmd := &cobra.Command{
		Use:     "list-project",
		Aliases: []string{"lp"},
		Short:   "List Project",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.accessToken = args[0]
			options.authClient = c.GetAuthServiceClient()
			options.projectCl = c.GetProjectServiceClient()
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runListProject(options)
		},
	}

	return cmd
}

func runListProject(options listingProjectOptions) error {
	req := &gateway.ListRequest{}

	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer " + options.accessToken})), time.Second*10)
	defer cancel()

	resp, err := options.projectCl.List(ctx, req)
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

	fmt.Println("List project result: ", resp.Projects)

	return nil
}
