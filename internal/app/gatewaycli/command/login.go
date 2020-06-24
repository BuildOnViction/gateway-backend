package command

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	gateway "github.com/anhntbk08/gateway/.gen/api/proto/gateway/v1"
	"github.com/spf13/cobra"
	"github.com/tomochain/tomochain/crypto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type loginOptions struct {
	address    string
	token      string
	privatekey string
	client     gateway.AuthServiceClient
}

// NewLoginCommand creates a new cobra.Command for adding a new item to the list.
func NewLoginCommand(c Context) *cobra.Command {
	options := loginOptions{}

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"login"},
		Short:   "Login",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.address = args[0]
			options.token = args[1]
			options.privatekey = args[2]
			options.client = c.GetAuthServiceClient()

			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			return runLogin(options)
		},
	}

	return cmd
}
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func runLogin(options loginOptions) error {
	privateKey, _ := crypto.HexToECDSA(options.privatekey)
	hexToken := []byte(options.token)
	signature, err := crypto.Sign(signHash(hexToken), privateKey)

	signature[64] += 27

	req := &gateway.AuthServiceLoginRequest{
		Address:   options.address,
		Token:     options.token,
		Signature: hex.EncodeToString(signature),
	}

	fmt.Println(
		options.address,
		options.token,
		hex.EncodeToString(signature),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := options.client.Login(ctx, req)
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

	fmt.Println("Login result: ", resp.AccessToken)

	return nil
}
