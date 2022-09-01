/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"
	"time"

	"github.com/spf13/cobra"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process messages in topic",
	Long:  "Process messages in topic",
	Run: func(cmd *cobra.Command, args []string) {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{*KafkaConfigInstance.Broker},
			Topic:   *KafkaConfigInstance.Topic,
		})
		defer reader.Close()

		cCluster := gocql.NewCluster(*cassandraHost)
		cCluster.Keyspace = *cassandraKeyspace
		session, err := cCluster.CreateSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()

		for {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			m, err := reader.ReadMessage(ctx)
			if err == context.DeadlineExceeded {
				fmt.Println("No new messages, exiting...")
				break
			} else if err != nil {
				panic(err)
			}

			var tweet Tweet
			if err := json.Unmarshal(m.Value, &tweet); err != nil {
				panic(err)
			}
			if err := session.Query("INSERT INTO tweets (id, created_at, text) VALUES (?, ?, ?)", string(m.Key), tweet.CreatedAtT(), tweet.Text).Exec(); err != nil {
				panic(err)
			}
		}
	},
}

type Tweet struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (t Tweet) CreatedAtT() time.Time {
	res, err := time.Parse("2006-01-02 15:04:05", t.CreatedAt)
	if err != nil {
		panic(err)
	}
	return res
}

var cassandraHost *string
var cassandraKeyspace *string

func init() {
	rootCmd.AddCommand(processCmd)

	cassandraHost = processCmd.Flags().String("cassandra-host", "localhost", "Cassandra host")
	cassandraKeyspace = processCmd.Flags().String("cassandra-keyspace", "gotest", "Cassandra keyspace")
}
