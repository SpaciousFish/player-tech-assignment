package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Useful variables when running tests
var authToken string = "Bearer abcd1234"
var host string = "http://localhost:5000"
var clientId string = "abcd"

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
	Test to check whether the create profile functionality is working.
*/
func TestCreateProfile(t *testing.T) {
	id := "1234"
	profile := `{"profile": {"applications": [{"applicationId": "test", "version": "1.0"}]}}`
	res, body := createProfile(id, profile)
	expected := `{"profile":{"applications":[{"applicationId":"test","version":"1.0"}]},"macAddress":"1234","clientId":"abcd"}`
	// remove newlines
	body = bytes.ReplaceAll(body, []byte("\n"), []byte(""))
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, expected, string(body))
	deleteProfile(id)
}

/**
	Test to check whether the get profile functionality is working.
*/
func TestGetProfile(t *testing.T) {
	id := "1234"
	profile := `{"profile": {"applications": [{"applicationId": "test", "version": "1.0"}]}}`
	createProfile(id, profile)
	res, body := executeRequest("GET", host+"/profiles/"+id, "")
	expected := `{"profile":{"applications":[{"applicationId":"test","version":"1.0"}]},"macAddress":"1234","clientId":"abcd"}`
	// remove newlines
	body = bytes.ReplaceAll(body, []byte("\n"), []byte(""))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, expected, string(body))
	deleteProfile(id)
}

/**
	Test to check whether the get all profiles functionality is working.
*/
func TestGetAllProfiles(t *testing.T) {
	id := "1234"
	profile := `{"profile": {"applications": [{"applicationId": "test", "version": "1.0"}]}}`
	createProfile(id, profile)
	res, body := executeRequest("GET", host+"/profiles", "")
	expected := `[{"profile":{"applications":[{"applicationId":"test","version":"1.0"}]},"macAddress":"1234","clientId":"abcd"}]`
	// remove newlines
	body = bytes.ReplaceAll(body, []byte("\n"), []byte(""))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, expected, string(body))
	deleteProfile(id)
}

/**
	Test to check whether the delete profile functionality is working.
*/
func TestDeleteProfile(t *testing.T) {
	id := "1234"
	profile := `{"profile": {"applications": [{"applicationId": "test", "version": "1.0"}]}}`
	createProfile(id, profile)
	res, body := deleteProfile(id)
	expected := `profile of client 1234 deleted`
	// remove newlines
	body = bytes.ReplaceAll(body, []byte("\n"), []byte(""))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, expected, string(body))
}

/**
	Test to check whether the update profile functionality is working.
*/
func TestUpdateProfile(t *testing.T) {
	id := "1234"
	profile := `{"profile": {"applications": [{"applicationId": "test", "version": "1.0"}]}}`
	createProfile(id, profile)
	profile = `{"profile": {"applications": [{"applicationId": "test", "version": "2.0"}]}}`
	res, body := executeRequest("PUT", host+"/profiles/"+id, profile)
	expected := `{"profile":{"applications":[{"applicationId":"test","version":"2.0"}]},"macAddress":"1234","clientId":"abcd"}`
	// remove newlines
	body = bytes.ReplaceAll(body, []byte("\n"), []byte(""))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, expected, string(body))
	deleteProfile(id)
}
