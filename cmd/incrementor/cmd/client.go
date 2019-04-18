package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"incrementor/internal/config"
	"incrementor/internal/models"
	"incrementor/internal/presentors/rpc"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "example call increments",
	Long:  "Call rpc server for incrementing numbers",
	Run:   clientMain,
}

// Append client cmd to root as a child
func init() {
	rootCmd.AddCommand(clientCmd)
}

// Client entrypoint
func clientMain(cmd *cobra.Command, args []string) {

	fmt.Println(models.MaxValue)

	var (
		err       error
		container *dig.Container
	)

	// Init Di container
	container = dig.New()

	// Append config.Config to di
	if err = container.Provide(config.NewConfig); err != nil {
		logrus.Fatalf("config: %s", err.Error())
	}

	// Append rpc.Server to di
	if err = container.Provide(rpc.NewIncrementorClient); err != nil {
		logrus.Fatalf("rpc.Client: %s", err.Error())
	}

	// Invoke server.Run from di
	if err = container.Invoke(func(client *rpc.IncrementorClient) error { return client.Run() }); err != nil {
		logrus.Fatalf("client.Run: %s", err.Error())
	}

}
