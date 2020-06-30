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

type getOneProjectOptions struct {
	accessToken string
	projectId   string
	authClient  gateway.AuthServiceClient
	projectCl   gateway.ProjectServiceClient
}

//  gets a new cobra.Command for adding a new item to the get.
func NewGetOneProjectCommand(c Context) *cobra.Command {
	options := getOneProjectOptions{}

	cmd := &cobra.Command{
		Use:     "get-project",
		Aliases: []string{"gp"},
		Short:   "Get Project",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.accessToken = args[0]
			options.projectId = args[1]
			options.authClient = c.GetAuthServiceClient()
			options.projectCl = c.GetProjectServiceClient()
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runGetOneProject(options)
		},
	}

	return cmd
}

func runGetOneProject(options getOneProjectOptions) error {
	req := &gateway.GetOneRequest{
		Id: options.projectId,
	}

	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer " + options.accessToken})), time.Second*10)
	defer cancel()

	resp, err := options.projectCl.GetOne(ctx, req)
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

	fmt.Printf("Get project result: %+v \n ", resp.Project)

	return nil
}
