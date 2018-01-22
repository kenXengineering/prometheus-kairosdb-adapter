package cmd

import (
	"github.com/chosenken/prometheus-kairosdb-adapter/pkg/adapter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Prints out received metrics from prometheus",
	Long:  `Prints out the received metrics from prometheus.  Usefull for testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		echoClient := adapter.NewEchoClient(&adapter.Options{
			PrintJson:  viper.GetBool("json"),
			ListenPort: viper.GetInt64("LISTEN_PORT"),
		})
		echoClient.Start()
	},
}

func init() {
	rootCmd.AddCommand(echoCmd)
	echoCmd.Flags().BoolP("json", "j", false, "Print messages in karisdb json format")
	viper.BindPFlag("json", echoCmd.Flags().Lookup("json"))
}
