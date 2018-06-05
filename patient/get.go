package patient

import (
	"encoding/json"
	"fmt"
	"github.com/callbinuvarghese/cassandra/PATIENT/cassandra"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

// Get -- handles GET request to /patients/ to fetch all Patients
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params (unused here)
func Get(w http.ResponseWriter, r *http.Request) {
	var patientList []Patient
	var count int
	m := map[string]interface{}{}

	fmt.Println("list-patients")
	query := "SELECT id,age,firstname,lastname,city,email FROM patient"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		fmt.Println("list-patients-iterating patient")
		patientList = append(patientList, Patient{
			ID:        m["id"].(gocql.UUID),
			Age:       m["age"].(int),
			FirstName: m["firstname"].(string),
			LastName:  m["lastname"].(string),
			Email:     m["email"].(string),
			City:      m["city"].(string),
		})
		m = map[string]interface{}{}
		count++
	}
	fmt.Println("list-patients: count ", count)

	json.NewEncoder(w).Encode(AllPatientsResponse{Patient: patientList})
}

// GetOne -- handles GET request to /patient/{patient_uuid} to fetch one patient
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func GetOne(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	var errs []string
	var found bool = false

	vars := mux.Vars(r)
	fmt.Println("list-patient:vars", vars)
	id := vars["patient_uuid"]

	fmt.Println("list-patient:", id)

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		fmt.Println("list-patient: id is not valid", id)
		errs = append(errs, err.Error())
	} else {
		fmt.Println("list-patient: retrieving wit valid id:", id)
		m := map[string]interface{}{}
		query := "SELECT id,age,firstname,lastname,city,email FROM patient WHERE id=? LIMIT 1"
		iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			fmt.Println("list-patient: iterate one")
			patient = Patient{
				ID:        m["id"].(gocql.UUID),
				Age:       m["age"].(int),
				FirstName: m["firstname"].(string),
				LastName:  m["lastname"].(string),
				Email:     m["email"].(string),
				City:      m["city"].(string),
			}
		}
		if !found {
			errs = append(errs, "Patient not found")
		}
	}

	if found {
		fmt.Println("list-patient: found patient", id)
		json.NewEncoder(w).Encode(GetPatientResponse{Patient: patient})
	} else {
		fmt.Println("list-patient: not found patient", id)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}

// Enrich -- turns an array of user UUIDs into a map of {uuid: "firstname lastname"}
// params:
// uuids - array of user UUIDs to fetch
// returns:
// a map[string]string of {uuid: "firstname lastname"}
func Enrich(uuids []gocql.UUID) map[string]string {
	if len(uuids) > 0 {
		fmt.Println("---\nfetching names", uuids)
		names := map[string]string{}
		m := map[string]interface{}{}

		query := "SELECT id,firstname,lastname FROM patients WHERE id IN ?"
		iterable := cassandra.Session.Query(query, uuids).Iter()
		for iterable.MapScan(m) {
			fmt.Println("m", m)
			patientID := m["id"].(gocql.UUID)
			fmt.Println("patientID", patientID.String())
			names[patientID.String()] = fmt.Sprintf("%s %s", m["firstname"].(string), m["lastname"].(string))
			m = map[string]interface{}{}
		}
		fmt.Println("names", names)
		return names
	}
	return map[string]string{}
}
