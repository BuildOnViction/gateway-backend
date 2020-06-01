package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anhntbk08/gateway/internal/app/gatewaycli"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gatewaycli",
		Short: "Gateway cli for authentication, list your projects settings, show list users, balance token ....",
	}

	gatewaycli.Configure(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
