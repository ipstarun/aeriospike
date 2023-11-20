// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const (
// 	dbHost = "localhost"
// 	dbName = "threejoin"
// )

// const numGoroutines = 1 // Number of concurrent goroutines

// var counter int // Global counter variable

// func main() {
// 	//Contexts are used for controlling and managing important aspects of reliable applications, 
// 	//such as cancellation and data sharing in concurrent programming

// 	//method to specify the connection details.
// 	clientOptions := options.Client().ApplyURI("mongodb://" + dbHost)
// 	client, err := mongo.Connect(context.Background(), clientOptions)

	

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err := client.Disconnect(context.Background()); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
		
// 	collection := client.Database(dbName).Collection("counter_collection")

// 	// Create the collection and insert one initial record with a counter value of 0
// 	_, err = collection.InsertOne(context.Background(), bson.M{"val": 0})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Channel to receive completion signals from goroutines
// 	done := make(chan bool)

// 	fmt.Println("dataBase connected successfully !!!!!!")

// 	// Start multiple goroutines to increment the counter concurrently
// 	for i := 0; i < numGoroutines; i++ {
// 		fmt.Println("starting time: " , time.Now())
// 		go func() {
// 			var timeSecond = time.Second
// 			startTime := time.Now()
// 			for time.Since(startTime) < timeSecond {
// 				_, err = collection.UpdateOne(
// 					context.Background(),
// 					bson.M{},
// 					bson.M{"$inc": bson.M{"val": 1}},
// 				)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				counter++
// 			}
// 			done <- true
// 			fmt.Printf("Incremented counter %d times.\n",  counter)
// 		}()
// 	}

// 	// Wait for all goroutines to complete
// 	for i := 0; i < numGoroutines; i++ {
// 		<-done
// 	}

// 	fmt.Println("Completion Time: ", time.Now())
// 	fmt.Printf("Final Counter Value: %d\n", counter)
// }



// /MONGODB ATLUS

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const (
// 	dbName = "threejoin" // Replace with your database name
// )

// const numGoroutines = 0 // Number of concurrent goroutines

// var counter int // Global counter variable

// func main() {
// 	// 	//Contexts are used for controlling and managing important aspects of reliable applications, 
// // 	//such as cancellation and data sharing in concurrent programming
// 	connectionString := "mongodb+srv://mongodb:mongodb@threejoin.b1nak1z.mongodb.net/?retryWrites=true&w=majority"

// 	clientOptions := options.Client().ApplyURI(connectionString)
// 	client, err := mongo.Connect(context.Background(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err := client.Disconnect(context.Background()); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	collection := client.Database(dbName).Collection("counter_collection")

// 	// Create the collection and insert one initial record with a counter value of 0
// 	_, err = collection.InsertOne(context.Background(), bson.M{"val": 0})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Channel to receive completion signals from goroutines
// 	done := make(chan bool)

// 	fmt.Println("Database connected successfully !!!!!!")

// 	// Start multiple goroutines to increment the counter concurrently
// 	for i := 0; i < numGoroutines; i++ {
// 		fmt.Println("starting time: ", time.Now())
// 		go func() {
// 			var timeSecond = time.Second
// 			startTime := time.Now()
// 			for time.Since(startTime) < timeSecond {
// 				_, err = collection.UpdateOne(
// 					context.Background(),
// 					bson.M{},
// 					bson.M{"$inc": bson.M{"val": 1}},
// 				)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				counter++
// 			}
// 			done <- true
// 			fmt.Printf("Incremented counter %d times.\n", counter)
// 		}()
// 	}

// 	// Wait for all goroutines to complete
// 	for i := 0; i < numGoroutines; i++ {
// 		<-done
// 	}

// 	fmt.Println("Completion Time: ", time.Now())
// 	fmt.Printf("Final Counter Value: %d\n", counter)
// }


package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "threejoin"
)

func main() {
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the counter table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS counter_table1 (
			id SERIAL PRIMARY KEY,
			value INT
		);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the counter value to 0 if it's not already set
	initCounterQuery := `
		INSERT INTO counter_table1 (value) VALUES (0)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err = db.Exec(initCounterQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Perform the benchmark
	startTime := time.Now()
	iterations := 0
	duration := 1 * time.Second

	for {
		// Update the counter value
		updateCounterQuery := `
			UPDATE counter_table1 SET value = value + 1 WHERE id = 1;
		`

		_, err := db.Exec(updateCounterQuery)
		if err != nil {
			log.Fatal(err)
		}

		iterations++

		// Check if 1 second has elapsed
		if time.Since(startTime) >= duration {
			break
		}
	}

	// Calculate and print the benchmark results
	elapsedTime := time.Since(startTime).Seconds()
	updatesPerSecond := float64(iterations) / elapsedTime

	fmt.Printf("Benchmark results:\n")
	fmt.Printf("  Elapsed Time: %.2f seconds\n", elapsedTime)
	fmt.Printf("  Updates per Second: %.2f\n", updatesPerSecond)
}


