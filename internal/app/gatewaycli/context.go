package gatewaycli

import (
	gateway "github.com/anhntbk08/gateway/.gen/api/proto/bridge/v1"
)

type context struct {
	client gateway.AuthServiceClient
}

func (c *context) GetAuthServiceClient() gateway.AuthServiceClient {
	return c.client
}
