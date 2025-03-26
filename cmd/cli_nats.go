package cmd

import (
	"log"

	"github.com/Phillezi/redis-vs-nats/pkg/bench"
	"github.com/Phillezi/redis-vs-nats/pkg/messaging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var natsCmd cobra.Command = cobra.Command{
	Use:   "nats",
	Short: "Benchmark NATS",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting NATS benchmark")
		log.Println("NATS host @ " + viper.GetString("nats-host"))
		connect := func() messaging.Broker {
			natsBroker, err := messaging.NewNATSBroker("nats://" + viper.GetString("nats-host"))
			if err != nil {
				log.Panic("could not connect to broker: " + err.Error())
				return nil // redundant
			}
			return natsBroker
		}

		bench.RunBenchmarks("NATS", connect)
	},
}

func init() {
	rootCmd.AddCommand(&natsCmd)
}
