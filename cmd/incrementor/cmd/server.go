package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"incrementor/internal/auth"
	"incrementor/internal/config"
	"incrementor/internal/database"
	"incrementor/internal/presentors/rpc"
	"incrementor/internal/repository"
	"incrementor/internal/services"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "serve increments",
	Long:  "Rpc server incrementing numbers for registererd & authneticated clients",
	Run:   serverMain,
}

// Append server cmd to root as a child
func init() {
	rootCmd.AddCommand(serverCmd)
}

// Server entrypoint
func serverMain(cmd *cobra.Command, args []string) {

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

	// Append sql.DB to di
	if err = container.Provide(database.NewConnection); err != nil {
		logrus.Fatalf("database: %s", err.Error())
	}

	// Append repositories to di
	if err = container.Provide(repository.NewClientRepositry); err != nil {
		logrus.Fatalf("repository.Client: %s", err.Error())
	}

	if err = container.Provide(repository.NewIncrementorRepositry); err != nil {
		logrus.Fatalf("repository.Incrementor: %s", err.Error())
	}

	// Append services to di
	if err = container.Provide(services.NewClientService); err != nil {
		logrus.Fatalf("services.Client: %s", err.Error())
	}

	if err = container.Provide(auth.NewJwt); err != nil {
		logrus.Fatalf("auth.Jwt: %s", err.Error())
	}
	if err = container.Provide(services.NewIncrementorService); err != nil {
		logrus.Fatalf("services.Incrementor: %s", err.Error())
	}

	if err = container.Provide(services.NewHealthChek); err != nil {
		logrus.Fatalf("services.Incrementor: %s", err.Error())
	}

	// Append rpc.Server to di
	if err = container.Provide(rpc.NewServer); err != nil {
		logrus.Fatalf("rc.server: %s", err.Error())
	}

	// Invoke server.Run from di
	if err = container.Invoke(func(server *rpc.Server) error { return server.Run() }); err != nil {
		logrus.Fatalf("server.Run: %s", err.Error())
	}
}
