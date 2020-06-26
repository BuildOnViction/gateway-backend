package command

import (
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
)

type issueAddressOptions struct {
	apikey     string
	authClient gateway.AuthServiceClient
	projectCl  gateway.ProjectServiceClient
}
