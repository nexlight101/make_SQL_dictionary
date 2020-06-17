package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// DData data structure for import
type DData map[string][]string

//Dict struct maps DData map to mongodb struct
type Dict struct {
	Word    string   `json:"word"`
	Explain []string `json:"explain"`
}

// Dictionary slice for dictionaries
type Dictionary []Dict

// port map to struct
func port(d map[string][]string) []Dict {
	record := Dict{}
	db := []Dict{}
	for key, value := range d {
		record.Word = key
		record.Explain = value
		db = append(db, record)
		record = Dict{} //empty record
	}
	return db
}

// createDB translates imported dictionary to
func createDB(d []Dict) {
	// Start sql client
	fmt.Println("Connecting to mySQL")
	// Open up our database connection.
	connString := "dictionaryuser:Password10@tcp(127.0.0.1:3306)/wordsdb?charset=utf8mb4"
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Cannot open connection: %v\n", err)
	}
	fmt.Println("opened connection")
	// defer the close till after the main function has finished
	defer db.Close()
	// test the link
	err = db.Ping()
	if err != nil {
		log.Fatalf("Not connected to mySQL: %v\n", err)
	}
	fmt.Println("Connected to mySQL")
	// Execute the query
	// results, err := db.Query("SELECT * FROM dictionary")
	// if err != nil {
	// 	log.Fatalf("Cannot execute query: %v\n", err)
	// }
	// for results.Next() {
	// 	record := Dict{}
	// 	// for each row, scan the result into our tag composite object
	// 	err = results.Scan(&record.Word, &record.Explain)
	// 	if err != nil {
	// 		panic(err.Error()) // proper error handling instead of panic in your app
	// 	}

	// }

	// Insert the data into wordsdb

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
	numV := "?, ?, ?, ?, ?, ?, ?"
	fields := "word,explain1,explain2,explain3,explain4,explain5,explain6"
	query := fmt.Sprintf("INSERT INTO dictionary(%s) VALUES(%s)", fields, numV)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Cannot prepare record into database: %v\n", err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	args := []string{}
	for _, v := range d {
		// build the query
		for _, item := range v.Explain {
			args = append(args, item)
		}
		for i := len(args); i < 6; i++ {
			args = append(args, "")
		}
		// fmt.Printf("record %d: %v\n", i, v)
		// fmt.Println(args)
		// execute the SQL query
		if _, err := stmt.Exec(v.Word, args[0], args[1], args[2], args[3], args[4], args[5]); err != nil {
			log.Fatalf("Cannot insert record into database: %v\n", err)
		}
		args = args[:0] //empty args to load new record
	}
}
func read(fName string, d DData) {
	file, err := ioutil.ReadFile(fName)
	if err != nil {
		log.Fatalf("Could not open file: %v\n", err)
	}

	jErr := json.Unmarshal(file, &d)
	if jErr != nil {
		log.Fatalf("Could not unmarshal the file: %v\n", err)
	}
}

func main() {
	dData := DData{}
	fName := "data/data.json"
	read(fName, dData)
	portedDB := port(dData)
	// for _, v := range portedDB {
	// 	fmt.Println(v)
	// }
	createDB(portedDB)
	fmt.Println("Database created")
}
