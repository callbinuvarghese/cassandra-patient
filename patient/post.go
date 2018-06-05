package patient

import (
	"encoding/json"
	"fmt"
	"github.com/callbinuvarghese/cassandra/PATIENT/cassandra"
	"github.com/gocql/gocql"
	"net/http"
)

/*
cul -X POST -H 'Content-Type: application/x-www-form-urlencoded' -d 'firstname=Ian&lastname=Douglas&city=Boulder&email=ian@getstream.io&age=42' "http://localhost:8080/patients/new"

*/
// Post -- handles POST request to /users/new to create new user
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUUID gocql.UUID

	patient, errs := FormToPatient(r)

	var created bool = false
	if len(errs) == 0 {
		fmt.Println("creating a new Patient")
		gocqlUUID = gocql.TimeUUID()
		if err := cassandra.Session.Query(`
		INSERT INTO Patient (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
			gocqlUUID, patient.FirstName, patient.LastName, patient.Email, patient.City, patient.Age).Exec(); err != nil {
			errs = append(errs, err.Error())

		} else {
			created = true
		}
	}

	if created {
		fmt.Println("user_id", gocqlUUID)
		json.NewEncoder(w).Encode(NewPatientResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
