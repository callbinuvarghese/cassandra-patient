package users

import (
	"encoding/json"
	"fmt"
	"github.com/callbinuvarghese/cassandra/PATIENT/cassandra"
	"github.com/gocql/gocql"
	"net/http"
)

/*
curl -X POST -H 'Content-Type: application/x-www-form-urlencoded'
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: X-Requested-With" \
  -X OPTIONS --verbose \
 -d 'firstname=Ben&lastname=Varghese&city=Cumming&email=ben.varghese@accenture.com&age=42' "http://localhost:8080/users"
{"id":"d966178c-6745-11e8-90ac-6a0001865cd0"}

*/
// Post -- handles POST request to /users/new to create new user
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUUID gocql.UUID

	user, errs := FormToUser(r)

	var created bool = false
	if len(errs) == 0 {
		fmt.Println("creating a new user")
		gocqlUUID = gocql.TimeUUID()
		if err := cassandra.Session.Query(`
		INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
			gocqlUUID, user.FirstName, user.LastName, user.Email, user.City, user.Age).Exec(); err != nil {
			errs = append(errs, err.Error())
			fmt.Println("Error creating a new user")
		} else {
			created = true
			fmt.Println("created a new user")
		}
	}

	if created {
		fmt.Println("user_id", gocqlUUID)
		json.NewEncoder(w).Encode(NewUserResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
