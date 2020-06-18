package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type updatingProjectOptions struct {
	accessToken string
	projectID   string
	field       string
	value       string
	authClient  gateway.AuthServiceClient
	projectCl   gateway.ProjectServiceClient
}

// NewUpdateprojectCommand updates a new cobra.Command for adding a new item to the update.
func NewUpdateProjectCommand(c Context) *cobra.Command {
	options := updatingProjectOptions{}

	cmd := &cobra.Command{
		Use:     "update-project",
		Aliases: []string{"up"},
		Short:   "Update Project",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.accessToken = args[0]
			options.projectID = args[1]
			options.field = args[2]
			options.value = args[3]
			options.authClient = c.GetAuthServiceClient()
			options.projectCl = c.GetProjectServiceClient()
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runUpdateProject(options)
		},
	}

	return cmd
}

func runUpdateProject(options updatingProjectOptions) error {
	req := gateway.UpdateRequest{
		Id: options.projectID,
	}

	if strings.ToLower(options.field) == "name" {
		req.Name = options.value
	}

	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer " + options.accessToken})), time.Second*10)
	defer cancel()

	resp, err := options.projectCl.Update(ctx, &req)
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

	fmt.Println("Update project result: ", resp.Success)

	return nil
}
