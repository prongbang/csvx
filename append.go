package csvx

import (
	"encoding/csv"
	"os"
)

// Append appends the given record to the end of the file specified by the filePath parameter. The record should be
// a slice of strings, where each string represents a field of the record. If the file does not exist, it will be created.
// If an error occurs during the operation, an error value will be returned.
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
