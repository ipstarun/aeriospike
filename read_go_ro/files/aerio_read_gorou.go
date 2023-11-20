package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	as "github.com/aerospike/aerospike-client-go"
)

func main() {
	// Connect to the Aerospike server
	client, err := as.NewClient("127.0.0.1", 3000)
	if err != nil {
		log.Fatalf("Error connecting to Aerospike: %v\n", err)
	}
	defer client.Close()

	// Set namespace, set name, and key
	namespace := "test"
	setName := "counter_set"
	key, err := as.NewKey(namespace, setName, "counter_key")
	if err != nil {
		log.Fatalf("Error creating key: %v\n", err)
	}

	// Set the initial value to 0 at the beginning of each run
	err = updateRecord(client, key, map[string]interface{}{"exampleBin": 0})
	if err != nil {
		log.Fatalf("Error setting initial value: %v\n", err)
	}

	// Run updates for one second
	updateOpsPerSecond := benchmarkUpdateOperations(client, key)
	fmt.Printf("Update operations per second: %d\n", updateOpsPerSecond)

	// Run reads for one second with goroutines
	readOpsPerSecond := benchmarkReadOperations(client, key)
	fmt.Printf("Read operations per second: %d\n", readOpsPerSecond)
}

func updateRecord(client *as.Client, key *as.Key, bins map[string]interface{}) error {
	writePolicy := as.NewWritePolicy(0, 0)
	err := client.Put(writePolicy, key, as.BinMap(bins))
	return err
}

func benchmarkUpdateOperations(client *as.Client, key *as.Key) int {
	// Start time for benchmarking
	startTime := time.Now()

	// Run updates for one second
	updateCount := 0
	for time.Since(startTime) < time.Second {
		// Increment the integer value in the bin
		err := updateRecord(client, key, map[string]interface{}{"exampleBin": 1})
		if err != nil {
			fmt.Printf("Error updating record: %v\n", err)
			return 0
		}
		updateCount++
	}

	// Print the final updated value
	fmt.Println("Update Starting time:", startTime)
	fmt.Println("Update Ending time:", time.Now())

	return updateCount
}

func benchmarkReadOperations(client *as.Client, key *as.Key) int {
	// Start time for benchmarking
	startTime := time.Now()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of goroutines (adjust according to your needs)
	numGoroutines := 10

	// Channel to collect results from goroutines
	resultCh := make(chan int, numGoroutines)

	// Run reads with goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			readCount := 0
			for time.Since(startTime) < time.Second {
				// Retrieve the value of the specified bin
				_, err := client.Get(nil, key)
				if err != nil {
					fmt.Printf("Error reading record: %v\n", err)
					return
				}
				readCount++
			}

			resultCh <- readCount
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the result channel to avoid deadlock
	close(resultCh)

	// Collect results from the channel
	totalReadCount := 0
	for result := range resultCh {
		totalReadCount += result
	}

	// Print the final updated value
	fmt.Println("Read Starting time:", startTime)
	fmt.Println("Read Ending time:", time.Now())

	return totalReadCount
}