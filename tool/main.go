package main

/*You need to create a production-ready tool that will automate the update of a thousand music players by using an API. You don't have to create the API.

Your tool will be used by different people using different operating systems. The most common ones will be Windows, MacOS and Linux.

The input is a .csv file containing, at the very minimum, MAC addresses of players to update, always in the first column.

### Example of a .csv file:
```
mac_addresses, id1, id2, id3
a1:bb:cc:dd:ee:ff, 1, 2, 3
a2:bb:cc:dd:ee:ff, 1, 2, 3
a3:bb:cc:dd:ee:ff, 1, 2, 3
a4:bb:cc:dd:ee:ff, 1, 2, 3
```

The `id1`, `id2` and `id3` fields aren't used in this assignment. The example is shown simply to demonstrate what the .csv file should look like.*/

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

func main() {
	// Open the file
	csvfile, err := os.Open("mac_addresses.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

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

			// Call the API to update the player
			url := "http://localhost:5000/profiles/" + record[0]
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(body)))
			if err != nil {
				fmt.Print(err.Error())
			}
			req.Header.Add("Authorization", authToken)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Print(err.Error())
			}
			defer res.Body.Close()
			body, readErr := io.ReadAll(res.Body)
			if readErr != nil {
				fmt.Print(readErr.Error())
			}
			fmt.Println(string(body))
		}
	}
}