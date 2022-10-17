package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var clients []Client = []Client{} 			// Our list of clients
var authToken string = "Bearer abcd1234"	// Bearer token used for requests

/**
	Structure for an application, contains an application id and a version.
*/
type Application struct {
	ApplicationId string `json:"applicationId"`
	Version       string `json:"version"`	
}

/**
	Structure for a profile, a profile contains many applications.
*/
type Profile struct {
	Applications []Application `json:"applications"`
}

/**
	Structure for a client, a client contains a profile, a mac addresss and a client id.
*/
type Client struct {
	Profile Profile `json:"profile"`
	MacAddress string `json:"macAddress"`
	ClientId string `json:"clientId"`
}

/**
	Function to create a profile
*/
func createProfile(w http.ResponseWriter, r *http.Request) {
	// check if the request is authorized
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid token supplied"))
		return
	}

	// check whether the client id is present
	if r.Header.Get("x-client-id") == "" {
		w.WriteHeader(400)
		w.Write([]byte("clientId not supplied"))
		return
	}

	// check if the client already exists
	params := mux.Vars(r)
	for _, item := range clients {
		if item.ClientId == r.Header.Get("x-client-id") || item.MacAddress == params["mac_address"] {
			w.WriteHeader(409)
			w.Write([]byte("client already exists"))
			return
		}
	}
	// Create a profile with a clientId and mac address
	var client Client
	_ = json.NewDecoder(r.Body).Decode(&client)
	client.ClientId = r.Header.Get("x-client-id")
	client.MacAddress = params["mac_address"]
	clients = append(clients, client)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(client)
}

/**
	Function to get a specific profile with a client id and mac address
*/
func getProfile(w http.ResponseWriter, r *http.Request) {
	// check if the request is authorized
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid clientId or token supplied"))
		return
	}

	// check whether the client id is present
	if r.Header.Get("x-client-id") == "" {
		w.WriteHeader(400)
		w.Write([]byte("clientId not supplied"))
		return
	}

	// look for the client
	params := mux.Vars(r)
	for _, item := range clients {
		if item.ClientId == r.Header.Get("x-client-id") && item.MacAddress == params["mac_address"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	// client not found, return error
	w.WriteHeader(404)
	w.Write([]byte("profile of client " + params["mac_address"] + " does not exist"))
}

/**
	Function to get all profiles in the application
*/
func getAllProfiles(w http.ResponseWriter, r *http.Request) {
	// check if the request is authorized
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid token supplied"))
		return
	}
	json.NewEncoder(w).Encode(clients)
}

/**
	Function to delete a profile given a client id and mac_address
*/
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	// check if the request is authorized
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid clientId or token supplied"))
		return
	}

	// check whether the client id is present
	if r.Header.Get("x-client-id") == "" {
		w.WriteHeader(400)
		w.Write([]byte("clientId not supplied"))
		return
	}

	// look for the client
	params := mux.Vars(r)
	for index, item := range clients {
		if item.ClientId == r.Header.Get("x-client-id") && item.MacAddress == params["mac_address"] {
			clients = append(clients[:index], clients[index+1:]...) // delete from the clients
			w.WriteHeader(200)
			w.Write([]byte("profile of client " + params["mac_address"] + " deleted"))
			return
		}
	}

	// client not found, return error
	w.WriteHeader(404)
	w.Write([]byte("profile of client " + params["mac_address"] + " does not exist"))
}

/**
	Function to update a profile given a client id and mac address
*/
func updateProfile(w http.ResponseWriter, r *http.Request) {
	// check if the request is authorized
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid clientId or token supplied"))
		return
	}

	// check whether the client id is present
	if r.Header.Get("x-client-id") == "" {
		w.WriteHeader(400)
		w.Write([]byte("clientId not supplied"))
		return
	}

	// look for the client
	params := mux.Vars(r)
	for index, item := range clients {
		if item.ClientId == r.Header.Get("x-client-id") && item.MacAddress == params["mac_address"] {
			var client Client
			_ = json.NewDecoder(r.Body).Decode(&client)

			// check whether applications are present
			if len(client.Profile.Applications) == 0 {
				w.WriteHeader(409)
				w.Write([]byte("child \"profile\" fails because [child \"applications\" fails because [\"applications\" is required]]"))
				return
			}

			// update the client
			clients = append(clients[:index], clients[index+1:]...)
			client.ClientId = r.Header.Get("x-client-id")
			client.MacAddress = params["mac_address"]
			clients = append(clients, client)
			json.NewEncoder(w).Encode(client)
			return
		}
	}

	// client not found, return error
	w.WriteHeader(404)
	w.Write([]byte("profile of client " + params["mac_address"] + " does not exist"))
}

/**
	This function is the main entry point of the program. It creates a router that maps each
	method above with a request that contains a method (post, get, put, delete) and a path.
	The server listens on port 5000 of the local computer.
*/
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/profiles/{mac_address}", createProfile).Methods("POST")
	router.HandleFunc("/profiles/{mac_address}", getProfile).Methods("GET")
	router.HandleFunc("/profiles", getAllProfiles).Methods("GET")
	router.HandleFunc("/profiles/{mac_address}", deleteProfile).Methods("DELETE")
	router.HandleFunc("/profiles/{mac_address}", updateProfile).Methods("PUT")
	http.ListenAndServe(":5000", router)
}