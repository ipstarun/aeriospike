package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbHost = "localhost"
	dbName = "threejoin"
)

const numGoroutines = 4 // Number of concurrent goroutines

var counter int // Global counter variable

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://" + dbHost)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database(dbName).Collection("counter_table")

	// Create the collection and insert one initial record with a counter value of 0
	_, err = collection.InsertOne(context.Background(), bson.M{"val": 0})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected successfully !!!!!!")

	// Test read operation within 1 second
	testReadOperation(collection)

	// Test insert operation within 1 second
	testInsertOperation(collection)

	fmt.Println("starting time is " , time.Now())

	// Start multiple goroutines to increment the counter concurrently
	done := make(chan bool)
	for i := 0; i < numGoroutines; i++ {
		
		go func() {
			var timeSecond = time.Second
			startTime := time.Now()
			for time.Since(startTime) < timeSecond {
				_, err = collection.UpdateOne(
					context.Background(),
					bson.M{},
					bson.M{"$inc": bson.M{"val": 1}},
				)
				if err != nil {
					log.Fatal(err)
				}
				counter++
			}
			done <- true
			fmt.Printf("Incremented counter %d times.\n", counter)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	fmt.Println("Completion Time: ", time.Now())
	fmt.Printf("Final Counter Value: %d\n", counter)
}

func testReadOperation(collection *mongo.Collection) {
	// Function to test read operation in 1 second and count the number of reads
	var timeSecond = time.Second
	startTime := time.Now()
	readCounter := 0
	for time.Since(startTime) < timeSecond {
		// Perform read operation (you can modify this part based on your actual read logic)
		// In this example, we're reading the document with an empty filter
		_, err := collection.FindOne(context.Background(), bson.M{}).DecodeBytes()
		if err != nil {
			log.Fatal(err)
		}
		readCounter++
	}
	fmt.Printf("Read operation tested successfully within 1 second. Read count: %d\n", readCounter)
}

func testInsertOperation(collection *mongo.Collection) {
	// Function to test insert operation in 1 second and count the number of inserts
	var timeSecond = time.Second
	startTime := time.Now()
	insertCounter := 0
	for time.Since(startTime) < timeSecond {
		_, err := collection.InsertOne(context.Background(), bson.M{"val": 1})
		if err != nil {
			log.Fatal(err)
		}
		insertCounter++
	}
	fmt.Printf("Insert operation tested successfully within 1 second. Insert count: %d\n", insertCounter)
}
