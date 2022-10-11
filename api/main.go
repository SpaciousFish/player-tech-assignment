package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var clients []Client = []Client{}
var authToken string = "Bearer abcd1234"

type Application struct {
	ApplicationId string `json:"applicationId"`
	Version       string `json:"version"`	
}

type Profile struct {
	Applications []Application `json:"applications"`
}

type Client struct {
	Profile Profile `json:"profile"`
	ClientId string `json:"clientId"`
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	// check if the request is authorized
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid token supplied"))
		return
	}
	// check if the client already exists
	params := mux.Vars(r)
	for _, item := range clients {
		if item.ClientId == params["clientId"] {
			w.WriteHeader(409)
			w.Write([]byte("client already exists"))
			return
		}
	}
	// Create a profile with a clientId
	var client Client
	_ = json.NewDecoder(r.Body).Decode(&client)
	// clientId is a mac address
	client.ClientId = params["clientId"]
	clients = append(clients, client)
	json.NewEncoder(w).Encode(client)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid clientId or token supplied"))
		return
	}
	params := mux.Vars(r)
	for _, item := range clients {
		if item.ClientId == params["clientId"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(404)
	w.Write([]byte("profile of client " + params["clientId"] + " does not exist"))
}

func getAllProfiles(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid token supplied"))
		return
	}
	json.NewEncoder(w).Encode(clients)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid clientId or token supplied"))
		return
	}
	params := mux.Vars(r)
	for index, item := range clients {
		if item.ClientId == params["clientId"] {
			clients = append(clients[:index], clients[index+1:]...)
			w.WriteHeader(200)
			w.Write([]byte("profile of client " + params["clientId"] + " deleted"))
			return
		}
	}
	w.WriteHeader(404)
	w.Write([]byte("profile of client " + params["clientId"] + " does not exist"))
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != authToken {
		w.WriteHeader(401)
		w.Write([]byte("invalid clientId or token supplied"))
		return
	}
	params := mux.Vars(r)
	for index, item := range clients {
		if item.ClientId == params["clientId"] {
			clients = append(clients[:index], clients[index+1:]...)
			var client Client
			_ = json.NewDecoder(r.Body).Decode(&client)
			client.ClientId = params["clientId"]
			clients = append(clients, client)
			json.NewEncoder(w).Encode(client)
			return
		}
	}
	w.WriteHeader(404)
	w.Write([]byte("profile of client " + params["clientId"] + " does not exist"))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/profiles/{clientId}", createProfile).Methods("POST")
	router.HandleFunc("/profiles/{clientId}", getProfile).Methods("GET")
	router.HandleFunc("/profiles", getAllProfiles).Methods("GET")
	router.HandleFunc("/profiles/{clientId}", deleteProfile).Methods("DELETE")
	router.HandleFunc("/profiles/{clientId}", updateProfile).Methods("PUT")
	http.ListenAndServe(":5000", router)
}