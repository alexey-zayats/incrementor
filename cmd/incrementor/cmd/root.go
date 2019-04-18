package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"incrementor/internal/config"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "incrementor",
	Short: "incrementor",
	Long:  "Provide incrementor rpc server for clients",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logrus.Fatal(err)
		}
	},
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

// Show cmd usage & exit
func usage(cmd *cobra.Command) {
	if err := cmd.Help(); err != nil {
		logrus.Fatal(err)
	}
	os.Exit(0)
}

// ConfigParam holds information about configuration options and mapping betwrrn Cobra cmd params and viper
type ConfigParam struct {
	Name      string
	Value     string
	Usage     string
	ViperBind string
}

// Init persistent flags, tie cobra and viper
func init() {

	cfgParams := []ConfigParam{
		{Name: "sql-driver", Value: "postgres", Usage: "sql driver", ViperBind: "sql.driver"},
		{Name: "sql-database", Value: "incrementor", Usage: "sql database", ViperBind: "sql.database"},
		{Name: "sql-username", Value: "incrementor", Usage: "sql username", ViperBind: "sql.username"},
		{Name: "sql-password", Value: "incrementor", Usage: "sql password", ViperBind: "sql.password"},
		{Name: "sql-hostname", Value: "localhost", Usage: "sql hostname", ViperBind: "sql.hostname"},
		{Name: "sql-port", Value: "5432", Usage: "sql port", ViperBind: "sql.port"},
		{Name: "sql-sslmode", Value: "disable", Usage: "sql ssl mode", ViperBind: "sql.sslmode"},

		{Name: "log-level", Value: "info", Usage: "log level", ViperBind: "log.level"},
		{Name: "log-target", Value: "postgres", Usage: "sql driver", ViperBind: "log.targert"},

		{Name: "server-listen-network", Value: "tcp", Usage: "server listen protocol", ViperBind: "server.listen.network"},
		{Name: "server-listen-address", Value: "127.0.0.1:9876", Usage: "server listen address", ViperBind: "server.listen.address"},
		{Name: "server-tls-certfile", Value: "config/tls/127.0.0.1.crt", Usage: "server listen address", ViperBind: "server.tls.certfile"},
		{Name: "server-tls-keyfile", Value: "config/tls/127.0.0.1.key", Usage: "server listen address", ViperBind: "server.tls.keyfile"},

		{Name: "client-auth-username", Value: "client", Usage: "client username to using incementor service", ViperBind: "client.auth.username"},
		{Name: "client-auth-password", Value: "secret-password", Usage: "client password to using incementor service", ViperBind: "client.auth.password"},
		{Name: "client-dial-address", Value: "127.0.0.1:9876", Usage: "client dial address to using incementor service", ViperBind: "client.dial.address"},
		{Name: "client-tls-certfile", Value: "config/tls/127.0.0.1.crt", Usage: "server listen address", ViperBind: "client.tls.certfile"},

		{Name: "incrementor-minvalue", Value: "0", Usage: "incrementor default minimum value", ViperBind: "incrementor.minvalue"},
		{Name: "incrementor-maxvalue", Value: "2147483647", Usage: "incrementor default maximum value (int32)", ViperBind: "incrementor.maxvalue"},
		{Name: "incrementor-incrementby", Value: "1", Usage: "incrementor default step", ViperBind: "incrementor.incrementby"},

		{Name: "jwt-secret", Value: "jwt-secret-key", Usage: "Secret for signing Jwt token", ViperBind: "jwt.secret"},
		{Name: "jwt-duration", Value: "1d", Usage: "Jwt token expired after", ViperBind: "jwt.duration"},
	}

	for _, p := range cfgParams {
		rootCmd.PersistentFlags().String(p.Name, p.Value, p.Usage)
		if err := viper.BindPFlag(p.ViperBind, rootCmd.PersistentFlags().Lookup(p.Name)); err != nil {
			usage(rootCmd)
		}
	}

	rootCmd.PersistentFlags().StringVar(&config.ConfigFile, "config", "./config/incrementor.yaml.dist", "Config file")

	cobra.OnInitialize(config.Init)
}
