package pkg

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// ReadCSV opens csv file defined in path and returns parsed strings.
func ReadCSV(path string) [][]string {

	var rawCSVData [][]string

	_, err := os.Stat(path)
	if err != nil {
		log.Println("Cannot stat", path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = -1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		rawCSVData = append(rawCSVData, strings.Split(record[0], ";"))
	}

	return rawCSVData
}
