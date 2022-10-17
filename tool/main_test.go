package main_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

var authToken string = "Bearer abcd1234"
var host string = "http://localhost:5000"

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
	url := host+"/profiles/" + mac_addr
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

func TestReadFile(t *testing.T) {
	r := readFile("mac_addresses.csv")
	assert.NotNil(t, r)
	record, err := r.Read()
	assert.Nil(t, err)
	assert.Equal(t, "mac_addresses", record[0])
}

func TestConnection(t *testing.T) {
	req, err := http.NewRequest("GET", host+"/profiles", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", authToken)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.NotNil(t, res.Body)
}

func TestSendRequestAndReturns401(t *testing.T) {
	var mac_addr string = "a1b2c3d4"
	var profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	var invalidClientId string = "a1"
	var invalidAuthToken string = "Bearer abcd12"
	// Call the API to update the player
	res, status := getResFromApi(mac_addr, profile, invalidAuthToken, invalidClientId)
	assert.Equal(t, "401 Unauthorized", status)
	assert.Equal(t, "invalid clientId or token supplied", res)
}

func TestSendRequestAndReturns404(t *testing.T) {
	var mac_addr string = "a1b2c3d4"
	var profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	var invalidClientId string = "a1"
	var validAuthToken string = "Bearer abcd1234"
	// Call the API to update the player
	res, status := getResFromApi(mac_addr, profile, validAuthToken, invalidClientId)
	assert.Equal(t, "404 Not Found", status)
	assert.Equal(t, "profile of client a1b2c3d4 does not exist", res)
}