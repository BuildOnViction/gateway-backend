package gatewaycli

import (
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
)

type context struct {
	client  gateway.AuthServiceClient
	project gateway.ProjectServiceClient
}

func (c *context) GetAuthServiceClient() gateway.AuthServiceClient {
	return c.client
}

func (c *context) GetProjectServiceClient() gateway.ProjectServiceClient {
	return c.project
}
