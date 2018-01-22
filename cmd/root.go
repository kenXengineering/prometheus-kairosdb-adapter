package cmd

import (
	"os"

	"fmt"

	"github.com/chosenken/prometheus-kairosdb-adapter/pkg/adapter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	log            = logrus.WithField("package", "main")
	cfgFile        string
	kairosDBServer string
	port           int64
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prometheus-kairosdb-adapter",
	Short: "Prometheus write adapter for KairosDB",
	Run:   run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&kairosDBServer, "kairosdb-url", "", "KairosDB URL")
	rootCmd.PersistentFlags().Int64VarP(&port, "listen-port", "p", 9201, "Listen Port")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable Debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("PKA")
	viper.AutomaticEnv()
	log.WithField("listenPort", port).Debug("Port value")
	if port != 0 {
		viper.Set("LISTEN_PORT", port)
	}
	if viper.GetInt64("LISTEN_PORT") == 0 {
		log.Error("Listen Port required.  Please set PKA_LISTEN_PORT or use the --listen-port / -p flag")
		os.Exit(1)
	}

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		log.Debug("Debug Enabled")
	}

}

// run will start the prometheus kairosdb adapter
func run(cmd *cobra.Command, args []string) {
	if len(kairosDBServer) != 0 {
		viper.Set("KAIROSDB_URL", kairosDBServer)
	}

	if len(viper.GetString("KAIROSDB_URL")) == 0 {
		log.Error("KairosDB URL required.  Please set PKA_KAIROSDB_URL or use the --kairosdb-url flag")
		os.Exit(1)
	}

	client, err := adapter.NewClient(&adapter.Options{
		ListenPort:  viper.GetInt64("LISTEN_PORT"),
		KairosDBURL: viper.GetString("KAIROSDB_URL"),
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	client.Start()
}
