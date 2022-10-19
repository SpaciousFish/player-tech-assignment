package main_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

// Useful variables
var authToken string = "Bearer abcd1234"
var host string = "http://localhost:5000"
var clientId string = "a1b2c3d4"

/**
	This function executes a request based on a keyword, a path, and a string body.
	Returns the result and the body.
*/
func executeRequest(keyword string, path string, body string) (*http.Response, []byte) {
	req, err := http.NewRequest(keyword, path, bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("x-client-id", clientId)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	res_body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print(readErr.Error())
	}
	return res, res_body
}

/**
	Helper function to create a profile based on a profile id and profile body.
	Returns the result and body of the request.
*/
func createProfile(id string, profile string) (*http.Response, []byte) {
	res, body := executeRequest("POST", host+"/profiles/"+id, profile)
	return res, body
}

/**
	Helper function to delete a profile based on a profile id and profile body.
	Returns the result and body of the request.
*/
func deleteProfile(id string) (*http.Response, []byte) {
	res, body := executeRequest("DELETE", host+"/profiles/"+id, "")
	return res, body
}

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
	Test whether the read file works correctly.
*/
func TestReadFile(t *testing.T) {
	r := readFile("mac_addresses.csv")
	assert.NotNil(t, r)
	record, err := r.Read()
	assert.Nil(t, err)
	assert.Equal(t, "mac_addresses", record[0])
}

/**
	Test whether the api is connected properly.
*/
func TestConnection(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:5000/profiles", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", authToken)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.NotNil(t, res.Body)
}

/**
	Test whether sending a update request with a valid bearer token returns a 200.
*/
func TestRequestAndReturns200(t *testing.T) {
	var mac_addr string = "a1"
	var profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	var validClientId string = clientId
	var validAuthToken string = authToken
	createProfile(mac_addr, profile)
	var new_profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.11"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	// Call the API to update the player
	_, status := getResFromApi(mac_addr, new_profile, validAuthToken, validClientId)
	assert.Equal(t, "200 OK", status)
	deleteProfile(mac_addr)
}

/**
	Test whether sending a request with an invalid bearer token returns a 401 error.
*/
func TestSendRequestAndReturns401(t *testing.T) {
	var mac_addr string = "a1"
	var profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	var invalidClientId string = "a1"
	var invalidAuthToken string = "Bearer abcd12"
	// Call the API to update the player
	res, status := getResFromApi(mac_addr, profile, invalidAuthToken, invalidClientId)
	assert.Equal(t, "401 Unauthorized", status)
	assert.Equal(t, "invalid clientId or token supplied", res)
}

/**
	Test whether sending a request with a client that does not exist returns a 404 error.
*/
func TestSendRequestAndReturns404(t *testing.T) {
	var invalid_mac_addr string = "a1"
	var profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	var invalidClientId string = "a1"
	var validAuthToken string = authToken
	// Call the API to update the player
	res, status := getResFromApi(invalid_mac_addr, profile, validAuthToken, invalidClientId)
	assert.Equal(t, "404 Not Found", status)
	assert.Equal(t, "profile of client " + invalid_mac_addr + " does not exist", res)
}

func TestSendRequestAndReturns409(t *testing.T) {
	var mac_addr string = "a1"
	var profile string = `{"profile":{"applications":[{"applicationId":"music_app","version":"v1.4.10"},{"applicationId":"diagnostic_app","version":"v1.2.6"},{"applicationId":"settings_app","version":"v1.1.5"}]}}`
	var validClientId string = clientId
	var validAuthToken string = authToken
	createProfile(mac_addr, profile)
	// Call the API to update the player
	var invalidProfile string = `{"profile":""}`
	res, status := getResFromApi(mac_addr, invalidProfile, validAuthToken, validClientId)
	assert.Equal(t, "409 Conflict", status)
	assert.Equal(t, "child \"profile\" fails because [child \"applications\" fails because [\"applications\" is required]]", res)
	deleteProfile(mac_addr)
}