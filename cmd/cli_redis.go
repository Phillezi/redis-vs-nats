package cmd

import (
	"log"

	"github.com/Phillezi/redis-vs-nats/pkg/bench"
	"github.com/Phillezi/redis-vs-nats/pkg/messaging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var redisCmd cobra.Command = cobra.Command{
	Use:   "redis",
	Short: "Benchmark Redis",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Redis benchmark")
		log.Println("Redis host @ " + viper.GetString("redis-host"))
		connect := func() messaging.Broker {
			return messaging.NewRedisBroker(viper.GetString("redis-host"))
		}
		bench.RunBenchmarks("Redis", connect)
	},
}

func init() {
	rootCmd.AddCommand(&redisCmd)
}
