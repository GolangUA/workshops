package main

import (
	"fmt"
	"log"

	"github.com/awnzl/workshops/cli-decoder/internal/filehandler"
	"github.com/awnzl/workshops/cli-decoder/internal/processor"
	"github.com/spf13/cobra"
)

func main() {
	var json, xml bool
	var help func() error

	var app = &cobra.Command{
		Use:   "cli-decoder",
		Short: "cli-decoder saves json or xml impute string to file",
		Run: func(cmd *cobra.Command, args []string) {
			if (json && xml) || (!json && !xml) {
				_ = help()
				fmt.Println("cli-decoder accepts only 1 arg specified above")
				return
			}

			if json {
				h := filehandler.NewJSONFilehandler()
				defer h.Release()
				work(processor.NewJSONProcessor(h))
			}

			if xml {
				h := filehandler.NewXMLFilehandler()
				defer h.Release()
				work(processor.NewXMLProcessor(h))
			}
		},
	}

	app.Flags().BoolVar(&json, "json", false, "process json string input")
	app.Flags().BoolVar(&xml, "xml", false, "process xml string input")
	help = app.Help

	err := app.Execute()
	if err != nil {
		log.Println(err)
	}
}

func work(p processor.Processor) {
	for {
		p.Process()
	}
}
