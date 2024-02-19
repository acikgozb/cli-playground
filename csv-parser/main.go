package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute operation")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	var operation statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w:%d", ErrInvalidColumn, column)
	}

	// Validate and define the operation
	switch op {
	case "sum":
		operation = sum
	case "avg":
		operation = avg
	default:
		return fmt.Errorf("%w:%s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	// Create the channels to receive results or errors of operations
	filesCh := make(chan string)
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{}) // Empty struct does not allocate any memory.

	wg := sync.WaitGroup{}

	// The main goroutine, the worker queues will pick from the channel to produce output
	go func() {
		defer close(filesCh)
		for _, fname := range filenames {
			filesCh <- fname
		}
	}()

	// Create worker queues, since this CLI is CPU bound, the upper limit of goroutines is tied to the CPU number of the working environment
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for fname := range filesCh {
				f, err := os.Open(fname)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
				}

				data, err := csv2float(f, column)
				if err != nil {
					errCh <- err
				}

				if err := f.Close(); err != nil {
					errCh <- err
				}

				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, operation(consolidate))
			return err
		}
	}
}
