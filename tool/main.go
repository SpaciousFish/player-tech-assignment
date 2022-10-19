package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var authToken string = "Bearer abcd1234" // Should be the same as the api auth token
var clientId string = "abcd"		 	 // Should be the same as is passed to the api

/**
	This function reads in a csv file from a specific path and returns the file reader.
*/
func readFile(path string) *csv.Reader {
	// Open the file
	csvfile, err := os.Open(path)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	return r
}

/**
	This function gets the result from the API when a PUT request is sent with a specific mac address,
	body, authentication token and client id. The return value is the body and the status of the result.
	If an error occurs, the body will be empty and the status will be Error.
*/
func getResFromApi(mac_addr string, bodyIn string, auth string, client string) (result string, status string){
	// Create the request
	url := "http://localhost:5000/profiles/" + mac_addr
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(bodyIn)))
	if err != nil {
		return "", "Error"
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("x-client-id", client)

	// Call the API to update the player
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "Error"
	}
	defer res.Body.Close()

	// Read the body contents
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", "Error"
	}

	return string(body), res.Status
	
}

/**
	This is the main function that happens on launch.
*/
func main() {
	// Open the file
	r := readFile("mac_addresses.csv")

	var body string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] != "mac_addresses" {
			// Print the record
			fmt.Println("MAC address:", record[0])
			res, status := getResFromApi(record[0], body, authToken, clientId)
			if status == "200 OK" {
				fmt.Println("Success:", status)
				fmt.Println("Response:", res)
			} else {
				fmt.Println("Error", status)
				fmt.Println("Response:", res)
			}
		}
	}
}