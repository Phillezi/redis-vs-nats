package cmd

import (
	"log"

	"github.com/Phillezi/redis-vs-nats/pkg/bench"
	"github.com/Phillezi/redis-vs-nats/pkg/messaging"
	"github.com/spf13/cobra"
)

var monoCmd cobra.Command = cobra.Command{
	Use:   "mono",
	Short: "Benchmark Monolithic with channels",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Mono benchmark")
		connect := func() messaging.Broker {
			return messaging.NewChannelBroker()
		}
		bench.RunBenchmarks("Mono", connect)
	},
}

func init() {
	rootCmd.AddCommand(&monoCmd)
}
