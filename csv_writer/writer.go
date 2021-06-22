package csv_writer

import (
	"encoding/csv"
	"os"
)

func WriteToCSV(fileName string, data [][]string) error {
	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	defer w.Flush()

	for _, record := range data {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil

}
