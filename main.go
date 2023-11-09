

package main

import as "github.com/aerospike/aerospike-client-go"
import (
	"fmt"
)


func main(){

	// Create new write policy
	policy := as.NewWritePolicy(0,0)
	policy.SendKey = true


	// Create new read policy
	readpolicy := as.NewPolicy()
	readpolicy.SocketTimeout = 300


	client, err := as.NewClient("127.0.0.1", 3000)

	if err != nil {
		fmt.Println("err is  " ,err)
	}
	defer client.Close()
	
	fmt.Println("client is  " ,client)


	key, err := as.NewKey("sandbox", "ufodata", 5001)
	if err != nil {
		fmt.Println("err is  " ,err)
	}

	fmt.Println("keys: ", key)

	// Create bin with new value
	newPosted := as.NewBin("balance", 100)

	client.PutBins(policy, key, newPosted)


	key1, err := as.NewKey("sandbox", "ufodata", 5001)
	if err != nil {
		fmt.Println("err is  " ,err)
	}

	// Get bins 'report' and 'location'
	record, err := client.Get(readpolicy, key1, "balance")
	if err != nil {
		fmt.Println("err is  " ,err)
	}

	// Do something
	fmt.Printf("Record: %v", record.Bins)
	

}