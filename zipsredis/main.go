package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

const usage = `
usage:
	zips load|get <zip>
`

//zipRecord represents one zip code record
type zipRecord struct {
	Zip   string
	City  string
	State string
}

//loadZips loads the zip codes from the csvFilePath into the redisClient
func loadZips(redisClient *redis.Client, csvFilePath string) (int, error) {
	f, err := os.Open(csvFilePath)
	if err != nil {
		return 0, fmt.Errorf("error opening zips file: %v", err)
	}

	//create a new CSV reader, which can read and parse
	//a stream of CSV data, one line at a time
	reader := csv.NewReader(f)

	//the first record is really the column names,
	//which we don't need, so just read and discard them
	_, err = reader.Read()
	if err != nil {
		return 0, fmt.Errorf("error reading CSV field names: %v", err)
	}

	fmt.Printf("loading zips")
	nZips := 0
	//read lines until we reach the end of the file
	//the .Read() method will return io.EOF when
	//it reaches the end of the file
	for {
		//read the next record
		record, err := reader.Read()
		//if we reached the end of the file,
		//return the zipSlice and no error
		if err == io.EOF {
			return nZips, nil
		}
		//if we encountered some other error,
		//return it
		if err != nil {
			return nZips, fmt.Errorf("error loading zips from CSV: %v", err)
		}

		//create and populate a new *zip
		//this syntax accomplishes the same thing
		//as using new() followed by field assignments
		//but does all of that in one statement
		zip := &zipRecord{
			Zip:   record[0],
			City:  record[3],
			State: record[6],
		}

		//add to redis database using zip as key
		if err := setZip(redisClient, zip); err != nil {
			return nZips, fmt.Errorf("error setting zip: %v", err)
		}
		nZips++
		if nZips%1000 == 0 {
			fmt.Print(".")
		}
	}
}

//setZip sets an entry in the redis database using the zip code
//as the key and the JSON-encoded zipRecord struct as the value
func setZip(redisClient *redis.Client, zip *zipRecord) error {
	return nil
}

//getZip gets the requested zip code form the redis database
//and returns the JSON-decoded zipRecord struct
func getZip(redisClient *redis.Client, zip string) (*zipRecord, error) {
	return nil, nil
}

//main is the main entry point for this program
func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])
}
