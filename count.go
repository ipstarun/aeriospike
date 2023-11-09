// for String type value 
// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/aerospike/aerospike-client-go/v5"
// )

// func main() {
// 	// Connect to the Aerospike cluster
// 	client, err := aerospike.NewClient("127.0.0.1", 3000)
// 	if err != nil {
// 		fmt.Printf("Error connecting to Aerospike: %v\n", err)
// 		return
// 	}

// 	// Set namespace and set name
// 	namespace := "test"
// 	setName := "exampleSet"

// 	// Record key and bin name
// 	key, _ := aerospike.NewKey(namespace, setName, "recordKey")
	
// 	//Specifies the name of the bin (a container for data within a record) as "exampleBin".
// 	binName := "exampleBin"

// 	// Define a sample record
// 	//Creates a map representing a sample record with a single bin named "exampleBin"
// 	//and an initial value of "initialValue".
// 	record := aerospike.BinMap{
// 		binName: "initialValue",
// 	}

// 	// Start time for benchmarking
// 	startTime := time.Now()

// 	// Number of updates to perform in one second
// 	updateCount := 0

// 	fmt.Println("staring time:  " ,  startTime)
// 	// Run updates for one second
// 	for time.Since(startTime) < time.Second {
		
// 		// Update the record
// 		err := client.Put(nil, key, record)
// 		if err != nil {
// 			fmt.Printf("Error updating record: %v\n", err)
// 			return
// 		}

// 		// Increment update count
// 		updateCount++
// 	}

// 	// Print benchmark results
// 	fmt.Printf("Updates per second: %d\n", updateCount)
// 	fmt.Println("ending time:  " ,  time.Now())
// }


// for interger type values
package main

import (
	"fmt"
	"time"

	"github.com/aerospike/aerospike-client-go/v5"
)

func main() {
	// Connect to the Aerospike cluster
	client, err := aerospike.NewClient("127.0.0.1", 3000)
	if err != nil {
		fmt.Printf("Error connecting to Aerospike: %v\n", err)
		return
	}

	// Set namespace and set name
	namespace := "test"
	setName := "counter_table"

	// Record key and bin name
	//namespace: Represents the namespace in the Aerospike database where the record is located.
	//recordKey:-Represents the actual key value for the specific record within the set, often a unique identifier.
	key, _ := aerospike.NewKey(namespace, setName, "recordKey")

	// Specifies the name of the bin as "exampleBin".
	binName := "exampleBin"

	// Define a sample record with an integer value
	// Creates a map representing a sample record with a single bin named "exampleBin"
	// and an initial integer value of 0.
	record := aerospike.BinMap{
		binName: 0,
	}

	// Start time for benchmarking
	startTime := time.Now()

	// Number of updates to perform in one second
	updateCount := 0

	fmt.Println("starting time: ", startTime)
	// Run updates for one second
	for time.Since(startTime) < time.Second {
		// Increment the integer value in the bin
		record[binName] = record[binName].(int) + 1

		// Update the record
		// It takes nil as the policy (indicating default policy), 
		//the key (key) of the record to be updated, and the updated record (record).
		err := client.Put(nil, key, record)
		if err != nil {
			fmt.Printf("Error updating record: %v\n", err)
			return
		}

		// Increment update count
		updateCount++
	}

	// Print benchmark results
	fmt.Printf("Updates per second: %d\n", updateCount)
	fmt.Println("ending time: ", time.Now())
}




