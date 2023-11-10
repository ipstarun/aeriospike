

// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// )

// func main() {
// 	// Connect to the Redis server
// 	client := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379", // Redis server address
// 		DB:   0,                // Default DB
// 	})

// 	// Close the connection when the main function exits
// 	defer client.Close()

// 	// Key and field names
// 	key := "counter_key"
// 	field := "exampleField"

// 	// Set the initial value to 0 at the beginning of each run
// 	_, err := client.HSet(context.Background(), key, field, 0).Result()
// 	if err != nil {
// 		fmt.Printf("Error setting initial value: %v\n", err)
// 		return
// 	}

// 	// Start time for benchmarking
// 	startTime := time.Now()

// 	// Run updates for one second
// 	for time.Since(startTime) < time.Second {
// 		// Increment the integer value in the field
// 		// Using HINCRBY command to increment the field value in a hash
// 		_, err := client.HIncrBy(context.Background(), key, field, 1).Result()
// 		if err != nil {
// 			fmt.Printf("Error updating record: %v\n", err)
// 			return
// 		}
// 	}

// 	// Get the final updated value
// 	finalValue, err := client.HGet(context.Background(), key, field).Result()
// 	if err != nil {
// 		fmt.Printf("Error getting final value: %v\n", err)
// 		return
// 	}

// 	// Print the final updated value
// 	fmt.Printf("Final Updated value: %s\n", finalValue)
// 	fmt.Println("Starting time:", startTime)
// 	fmt.Println("Ending time:", time.Now())
// }

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

func updateCounter(client *redis.Client, key, field string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Run updates for one second
	startTime := time.Now()
	updateCount := 0

	for time.Since(startTime) < time.Second {
		// Increment the integer value in the field
		_, err := client.HIncrBy(context.Background(), key, field, 1).Result()
		if err != nil {
			fmt.Printf("Error updating record: %v\n", err)
			return
		}

		// Increment update count
		updateCount++
	}

	// Print benchmark results
	fmt.Printf("Updates per second: %d\n", updateCount)
	fmt.Println("Thread ending time:", time.Now())
}

func main() {
	// Connect to the Redis server
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Default DB
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

	// Number of goroutines (threads) to run concurrently
	numGoroutines := 16
	var wg sync.WaitGroup

	// Print initial time
	fmt.Println("Initial time:", time.Now())

	// Start multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go updateCounter(client, key, field, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Get the final updated value
	finalValue, err := client.HGet(context.Background(), key, field).Result()
	if err != nil {
		fmt.Printf("Error getting final value: %v\n", err)
		return
	}

	// Print the final updated value
	fmt.Printf("Final Updated value: %s\n", finalValue)
	fmt.Println("Main thread ending time:", time.Now())
}
