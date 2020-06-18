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

type deletingProjectOptions struct {
	accessToken string
	projectID   string
	field       string
	value       string
	authClient  gateway.AuthServiceClient
	projectCl   gateway.ProjectServiceClient
}

// NewDeleteprojectCommand deletes a new cobra.Command for adding a new item to the delete.
func NewDeleteProjectCommand(c Context) *cobra.Command {
	options := deletingProjectOptions{}

	cmd := &cobra.Command{
		Use:     "delete-project",
		Aliases: []string{"dp"},
		Short:   "Delete Project",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.accessToken = args[0]
			options.projectID = args[1]
			options.authClient = c.GetAuthServiceClient()
			options.projectCl = c.GetProjectServiceClient()
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runDeleteProject(options)
		},
	}

	return cmd
}

func runDeleteProject(options deletingProjectOptions) error {
	req := gateway.DeleteRequest{
		Id: options.projectID,
	}

	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer " + options.accessToken})), time.Second*10)
	defer cancel()

	resp, err := options.projectCl.Delete(ctx, &req)
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

	fmt.Println("Delete project result: ", resp.Success)

	return nil
}
