package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
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

const numGoroutines = 4 // Number of concurrent goroutines

var counter int // Global counter variable

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table and insert one initial record with a counter value of 0
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS counter_table (val INT)")
	if err != nil {
		log.Fatal(err)
	}

	// Create the table and insert one initial record with a counter value of 0
	_, err = db.Exec("INSERT INTO counter_table (val) VALUES (0) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatal(err)
	}

	// Channel to receive completion signals from goroutines
	done := make(chan bool)
	fmt.Println("Starting time ", time.Now())

	// Start multiple goroutines to increment the counter concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			var timeSecond = time.Second
			startTime := time.Now()
			for time.Since(startTime) < timeSecond {
				_, err = db.Exec("UPDATE counter_table SET val = val + 1")
				if err != nil {
					log.Fatal(err)
				}
				counter++
			}
			done <- true
			fmt.Printf("Goroutine %d completed. Incremented counter %d times.\n", i, counter)
		}()
	}

	// Benchmark read operations
	benchmarkReadOperations(db)

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	fmt.Println("Completion Time: ", time.Now())
	fmt.Printf("Final Counter Value: %d\n", counter)
}

func benchmarkReadOperations(db *sql.DB) {
	var wg sync.WaitGroup
	readCounter := 0

	// Start multiple goroutines to perform SELECT queries concurrently
	for i := 0; i < numGoroutines; i++ {
		 // Increment the WaitGroup counter to indicate a new goroutine is started
		wg.Add(1)
		go func() {
			defer wg.Done()
			var timeSecond = time.Second
			startTime := time.Now()
			for time.Since(startTime) < timeSecond {
				var val int
				err := db.QueryRow("SELECT val FROM counter_table").Scan(&val)
				if err != nil {
					log.Fatal(err)
				}
				readCounter++
			}
		}()
	}

	// Wait for all read goroutines to complete
	wg.Wait()
	fmt.Printf("Benchmark: Read Counter Value: %d\n", readCounter)
}
