package main_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var authToken string = "Bearer abcd1234"
var host string = "http://localhost:5000"
var clientId string = "abcd"

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

func createProfile(id string, profile string) (*http.Response, []byte) {
	res, body := executeRequest("POST", host+"/profiles/"+id, profile)
	return res, body
}

func deleteProfile(id string) (*http.Response, []byte) {
	res, body := executeRequest("DELETE", host+"/profiles/"+id, "")
	return res, body
}

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
