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

var authToken string = "Bearer abcd1234"
var clientId string = "a1b2c3d4"

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

func getResFromApi(mac_addr string, bodyIn string, auth string, client string) (result string, status string){
	// Call the API to update the player
	url := "http://localhost:5000/profiles/" + mac_addr
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(bodyIn)))
	if err != nil {
		return "", "Error"
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("x-client-id", client)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "Error"
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", "Error"
	}
	return string(body), res.Status
	
}

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

