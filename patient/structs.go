package patient

import (
	"github.com/gocql/gocql"
)

// Patient struct to hold profile data for our patient
type Patient struct {
	ID        gocql.UUID `json:"id"`
	FirstName string     `json:"firstname"`
	LastName  string     `json:"lastname"`
	Email     string     `json:"email"`
	Age       int        `json:"age"`
	City      string     `json:"city"`
}

// GetPatientResponse to form payload returning a single Patient struct
type GetPatientResponse struct {
	Patient Patient `json:"patient"`
}

// AllPatientsResponse to form payload of an array of Patient structs
type AllPatientsResponse struct {
	Patient []Patient `json:"patient"`
}

// NewPatientResponse builds a payload of new Patient resource ID
type NewPatientResponse struct {
	ID gocql.UUID `json:"id"`
}

// ErrorResponse returns an array of error strings if appropriate
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
