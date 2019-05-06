package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"incrementor/internal/config"
	"os"
	"strings"
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
		{Name: "debug", Value: "true", Usage: "show debug information", ViperBind: "Debug"},

		{Name: "sql-driver", Value: "postgres", Usage: "sql driver", ViperBind: "Sql.Driver"},
		{Name: "sql-database", Value: "incrementor", Usage: "sql database", ViperBind: "Sql.Database"},
		{Name: "sql-username", Value: "incrementor", Usage: "sql username", ViperBind: "Sql.Username"},
		{Name: "sql-password", Value: "incrementor", Usage: "sql password", ViperBind: "Sql.Password"},
		{Name: "sql-hostname", Value: "localhost", Usage: "sql hostname", ViperBind: "Sql.hostname"},
		{Name: "sql-port", Value: "5432", Usage: "sql port", ViperBind: "Sql.Port"},
		{Name: "sql-sslmode", Value: "disable", Usage: "sql ssl mode", ViperBind: "Sql.SslMode"},

		{Name: "log-level", Value: "info", Usage: "log level", ViperBind: "Log.Level"},
		{Name: "log-target", Value: "postgres", Usage: "sql driver", ViperBind: "Log.Targert"},

		{Name: "server-listen-network", Value: "tcp", Usage: "server listen protocol", ViperBind: "Server.Listen.Network"},
		{Name: "server-listen-address", Value: "127.0.0.1:9876", Usage: "server listen address", ViperBind: "Server.Listen.Address"},
		{Name: "server-tls-certfile", Value: "config/tls/127.0.0.1.crt", Usage: "server listen address", ViperBind: "Server.TLS.CertFile"},
		{Name: "server-tls-keyfile", Value: "config/tls/127.0.0.1.key", Usage: "server listen address", ViperBind: "Server.TLS.KeyFile"},

		{Name: "client-auth-username", Value: "client", Usage: "client username to using incementor service", ViperBind: "Client.Auth.Username"},
		{Name: "client-auth-password", Value: "secret-password", Usage: "client password to using incementor service", ViperBind: "Client.Auth.Password"},
		{Name: "client-dial-address", Value: "127.0.0.1:9876", Usage: "client dial address to using incementor service", ViperBind: "Client.Dial.Address"},
		{Name: "client-tls-certfile", Value: "config/tls/127.0.0.1.crt", Usage: "server listen address", ViperBind: "client.TLS.CertFile"},

		{Name: "incrementor-minvalue", Value: "0", Usage: "incrementor default minimum value", ViperBind: "Incrementor.MinValue"},
		{Name: "incrementor-maxvalue", Value: "2147483647", Usage: "incrementor default maximum value (int32)", ViperBind: "Incrementor.MaxValue"},
		{Name: "incrementor-incrementby", Value: "1", Usage: "incrementor default step", ViperBind: "Incrementor.IncrementBy"},

		{Name: "jwt-secret", Value: "jwt-secret-key", Usage: "Secret for signing Jwt token", ViperBind: "JWT.Secret"},
		{Name: "jwt-duration", Value: "1d", Usage: "Jwt token expired after", ViperBind: "JWT.Duration"},
	}

	viper.SetEnvPrefix(config.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, p := range cfgParams {
		rootCmd.PersistentFlags().String(p.Name, p.Value, p.Usage)
		if err := viper.BindPFlag(p.ViperBind, rootCmd.PersistentFlags().Lookup(p.Name)); err != nil {
			usage(rootCmd)
		}
		viper.SetDefault(p.ViperBind, p.Value)
	}

	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&config.ConfigFile, "config", "./config/incrementor.yaml", "Config file")

	cobra.OnInitialize(config.Init)
}
