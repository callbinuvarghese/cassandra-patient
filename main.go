package main

/*
https://getstream.io/blog/building-a-performant-api-using-go-and-cassandra/

echo "CREATE KEYSPACE test WITH replication = {'class': 'NetworkTopologyStrategy', 'Analytics': '3'}  AND durable_writes = true;" | cqlsh

use test;
CREATE TABLE patients (
id UUID,
firstname text,
lastname text,
age int,
email text,
city text,
PRIMARY KEY (id)
);
CREATE TABLE users (
id UUID,
firstname text,
lastname text,
age int,
email text,
city text,
PRIMARY KEY (id)
);

*/

import (
	"encoding/json"
	"fmt"
	"github.com/callbinuvarghese/cassandra/PATIENT/cassandra"
	"github.com/callbinuvarghese/cassandra/PATIENT/patient"
	"github.com/callbinuvarghese/cassandra/PATIENT/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

const AppPort = ":8080"

func main() {
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)

	router.HandleFunc("/patients", patient.Get)
	//if you need get and post for the same mapping, please use as below
	//router.HandleFunc("/patients", patient.Get).Methods("GET")
	//router.HandleFunc("/patients", patient.Post).Methods("POST")

	router.HandleFunc("/patients/new", patient.Post)
	router.HandleFunc("/patients/{patient_uuid}", patient.GetOne)

	router.HandleFunc("/users", users.Get).Methods("GET")
	router.HandleFunc("/users", users.Post).Methods("POST")
	router.HandleFunc("/users/{user_uuid}", users.GetOne)

	fmt.Println("Server listening" + AppPort)
	//log.Fatal(http.ListenAndServe(AppPort, router))
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	//for getting user agent in the log
	//loggedRouter := CombinedLoggingHandler(os.Stdout, router)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	//originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	//log.Fatal(http.ListenAndServe(AppPort,
	//	handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"}), handlers.AllowedOrigins([]string{"*"}))(loggedRouter)))

	log.Fatal(http.ListenAndServe(AppPort,
		handlers.CORS(headersOk, originsOk, methodsOk)(loggedRouter)))

}

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
