package processor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/awnzl/workshops/cli-decoder/internal/filehandler"
)

type Processor interface {
	Process()
}

// JSON Processor
type JSONProcessor struct {
	filehandler filehandler.Filehandler
}

func NewJSONProcessor(fh filehandler.Filehandler) *JSONProcessor {
	return &JSONProcessor{
		filehandler: fh,
	}
}

func (p *JSONProcessor) Process() {
	printGreenStr("Please enter a valid JSON string and than press Enter key.")

	bts, err := readStdin()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !json.Valid(bts) {
		fmt.Println("Error:", "json data isn't valid")
		return
	}

	printGreenStr("The data is validated.")

	err = p.filehandler.SaveBytes(bts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printGreenStr("The data is saved.")
}

// XML Processor
type XMLProcessor struct {
	filehandler filehandler.Filehandler
}

func NewXMLProcessor(fh filehandler.Filehandler) *XMLProcessor {
	return &XMLProcessor{
		filehandler: fh,
	}
}

func (p *XMLProcessor) Process() {
	printGreenStr("Please enter a valid JSON string and than press Enter key.")

	bts, err := readStdin()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	decoder := xml.NewDecoder(bytes.NewReader(bts))

	for {
		err := decoder.Decode(new(interface{}))
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error:", "xml data isn't valid")
			return
		}
	}

	printGreenStr("The data is validated.")

	err = p.filehandler.SaveBytes(bts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printGreenStr("The data is saved.")
}

func readStdin() ([]byte, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return scanner.Bytes(), nil
}

func printGreenStr(s string) {
	fmt.Printf("\033[32m%v\033[0m\n", s)
}
