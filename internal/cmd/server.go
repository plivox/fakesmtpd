package cmd

import (
	"fakesmtpd/internal/server"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagAddress     string
	flagDomain      string
	flagTLS         bool
	flagTLSInsecure bool
	flagCACert      string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start fake SMTP Server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.NewServer(globalConfig); err != nil {
			log.Fatal().Err(err).Msg("Server: cannot start")
		}
	},
}

func init() {
	serverCmd.PersistentFlags().
		StringVar(&flagAddress, "address", "127.0.0.1:2525", "Server listen address")
	serverCmd.PersistentFlags().
		StringVar(&flagDomain, "domain", "localhost", "Server domain")
	serverCmd.PersistentFlags().
		BoolVar(&flagTLS, "tls", false, "Enable incoming TLS connections")
	serverCmd.PersistentFlags().
		BoolVar(&flagTLSInsecure, "insecure", false, "Allow insecure TLS connections")
	serverCmd.PersistentFlags().
		StringVar(&flagCACert, "cacert", "", "CA certificate to verify peer against")

	rootCmd.AddCommand(serverCmd)
}
