package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func sum(data []float64) float64 {
	sum := 0.0

	for _, float := range data {
		sum += float
	}

	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

// Define a generic statistical function type
type statsFunc func(data []float64) float64

func csv2float(r io.Reader, column int) ([]float64, error) {
	// Create the CSV reader used to read in data from CSV files
	csvReader := csv.NewReader(r)

	// Adjusting for 0 based index
	column--

	// Read in all CSV data
	allData, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}

	// Convert string to float64 to make calculations on it
	var data []float64

	for i, row := range allData {
		if i == 0 {
			continue
		}

		// Checking number of columns in CSV file
		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		// Try to convert data read into a float number
		float, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, float)
	}

	return data, nil
}
