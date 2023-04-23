package csvx

import (
	"encoding/csv"
	"os"
)

func Append(filePath string, record []string) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		_, err = os.Create(filePath)
		if err != nil {
			return err
		}
	}

	// Open the CSV file for writing and appending
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Create a new CSV writer and pass in the opened file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the new row to the CSV file
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}
