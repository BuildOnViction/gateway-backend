package command

import (
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	"github.com/spf13/cobra"
)

// Context represents the application context.
type Context interface {
	GetAuthServiceClient() gateway.AuthServiceClient
	GetProjectServiceClient() gateway.ProjectServiceClient
}

// AddCommands adds all the commands from cli/command to the root command.
func AddCommands(cmd *cobra.Command, c Context) {
	cmd.AddCommand(
		NewRequestTokenCommand(c),
		NewLoginCommand(c),
		NewCreateProjectCommand(c),
		NewListProjectCommand(c),
		NewUpdateProjectCommand(c),
		NewDeleteProjectCommand(c),
		NewGetOneProjectCommand(c),
	)
}
