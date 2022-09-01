package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const batchSize = 500

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load CSV data to kafka",
	Long:  "Load CSV data to kafka",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(*InputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer = KafkaConfigInstance.NewWriter()
		defer writer.Close()

		var msgs []kafka.Message
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if len(headers) == 0 {
				parseHeader(scanner.Text())
				continue
			}
			key, value := parseLine(scanner.Text())
			jsonVal, err := json.Marshal(value)
			if err != nil {
				panic(err)
			}
			msgs = append(msgs, kafka.Message{
				Key:   []byte(key),
				Value: jsonVal,
			})

			if len(msgs) >= batchSize {
				err := writer.WriteMessages(context.Background(), msgs...)
				if err != nil {
					panic(err)
				}
				msgs = []kafka.Message{}
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		if len(msgs) >= batchSize {
			err := writer.WriteMessages(context.Background(), msgs...)
			if err != nil {
				panic(err)
			}
			msgs = []kafka.Message{}
		}
	},
}

var InputFile *string
var IDColumn *string
var Separator *string

var headers []string
var idColumnIndex int
var writer *kafka.Writer

func parseHeader(value string) {
	parts := strings.Split(value, *Separator)

	for i, part := range parts {
		if part == *IDColumn {
			idColumnIndex = i
		}
		headers = append(headers, part)
	}

	if idColumnIndex == -1 {
		panic("ID column not found")
	}
}

func parseLine(value string) (string, map[string]string) {
	parts := strings.Split(value, *Separator)
	id := parts[idColumnIndex]
	data := make(map[string]string)
	for i, part := range parts {
		if i != idColumnIndex {
			data[headers[i]] = part
		}
	}
	return id, data
}

func Write(key string, value map[string]string) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: jsonVal,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(loadCmd)

	InputFile = loadCmd.Flags().StringP("file", "f", "", "Input file")
	IDColumn = loadCmd.Flags().StringP("id-column", "i", "", "ID column name")
	Separator = loadCmd.Flags().StringP("separator", "s", "\t", "Separator")

	loadCmd.MarkFlagRequired("file")
	loadCmd.MarkFlagRequired("id-column")
}
