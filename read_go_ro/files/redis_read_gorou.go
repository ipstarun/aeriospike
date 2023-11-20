package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Connect to the Redis server
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Default DB This sets the database index to 0
	})

	// Close the connection when the main function exits
	defer client.Close()

	// Key and field names
	key := "counter_key"
	field := "exampleField"

	// Set the initial value to 0 at the beginning of each run
	_, err := client.HSet(context.Background(), key, field, 0).Result()
	if err != nil {
		fmt.Printf("Error setting initial value: %v\n", err)
		return
	}

	// Run updates for one second
	updateOpsPerSecond := benchmarkUpdateOperations(client, key, field)
	fmt.Printf("Update operations per second: %d\n", updateOpsPerSecond)

	// Run reads for one second using goroutines
	readOpsPerSecond := benchmarkReadOperationsWithGoroutines(client, key, field)
	fmt.Printf("Read operations per second: %d\n", readOpsPerSecond)
}

func benchmarkUpdateOperations(client *redis.Client, key, field string) int {
	// Start time for benchmarking
	startTime := time.Now()

	// Run updates for one second
	updateCount := 0
	for time.Since(startTime) < time.Second {
		// Increment the integer value in the field
		_, err := client.HIncrBy(context.Background(), key, field, 1).Result()

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

func benchmarkReadOperationsWithGoroutines(client *redis.Client, key, field string) int {
	// Start time for benchmarking
	startTime := time.Now()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Channel to receive results from goroutines
	resultChan := make(chan int)

	// Run reads using goroutines
	for i := 0; i < 4; i++ { // Adjust the number of goroutines as needed
		wg.Add(1)
		go func() {
			defer wg.Done()

			readCount := 0
			for time.Since(startTime) < time.Second {
				// Retrieve the value of the specified field in a Redis Hash
				_, err := client.HGet(context.Background(), key, field).Result()

				if err != nil {
					fmt.Printf("Error reading record: %v\n", err)
					return
				}

				readCount++
			}

			resultChan <- readCount
		}()
	}

	// Close the result channel when all goroutines finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from goroutines
	totalReads := 0
	for reads := range resultChan {
		totalReads += reads
	}

	// Print the final updated value
	fmt.Println("Read Starting time:", startTime)
	fmt.Println("Read Ending time:", time.Now())

	return totalReads
}
