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
	csvReader.ReuseRecord = true

	// Adjusting for 0 based index
	column--

	var data []float64

	// As we can see from the benchmarks, it is not healthy to store all the data into a slice.
	// Therefore, change the ReadAll() call to read record by record
	for i := 0; ; i++ {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}

		// Skip first row (which includes column names)
		if i == 0 {
			continue
		}

		if len(row) <= column {
			return nil, fmt.Errorf("%w: file has only %d columns", ErrInvalidColumn, len(row))
		}

		// Convert string to float64 to make calculations on it
		convertedValue, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w:%s", ErrNotNumber, err)
		}

		data = append(data, convertedValue)
	}

	return data, nil
}
