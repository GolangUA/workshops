package cmd

import (
	"github.com/segmentio/kafka-go"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "data-pipeline",
	Short: "data-pipeline workshop",
	Long:  "data-pipeline workshop",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type KafkaConfig struct {
	Broker *string
	Topic  *string
}

func (kc *KafkaConfig) NewWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(*kc.Broker),
		Topic:    *kc.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}

var KafkaConfigInstance *KafkaConfig

func init() {
	KafkaConfigInstance = &KafkaConfig{
		Broker: loadCmd.Flags().StringP("broker", "b", "localhost:9092", "kafka broker"),
		Topic:  loadCmd.Flags().StringP("topic", "t", "tweets", "topic name"),
	}
}
