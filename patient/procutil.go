package patient

import (
	"net/http"
	"strconv"
)

// FormToPatient -- fills a Patient struct with submitted form data
// params:
// r - request reader to fetch form data or url params (unused here)
// returns:
// User struct if successful
// array of strings of errors if any occur during processing
func FormToPatient(r *http.Request) (Patient, []string) {
	var patient Patient
	var errStr, ageStr string
	var errs []string
	var err error

	patient.FirstName, errStr = processFormField(r, "firstname")
	errs = appendError(errs, errStr)
	patient.LastName, errStr = processFormField(r, "lastname")
	errs = appendError(errs, errStr)
	patient.Email, errStr = processFormField(r, "email")
	errs = appendError(errs, errStr)
	patient.City, errStr = processFormField(r, "city")
	errs = appendError(errs, errStr)

	ageStr, errStr = processFormField(r, "age")
	if len(errStr) != 0 {
		errs = append(errs, errStr)
	} else {
		patient.Age, err = strconv.Atoi(ageStr)
		if err != nil {
			errs = append(errs, "Parameter 'age' not an integer")
		}
	}
	return patient, errs
}

func appendError(errs []string, errStr string) []string {
	if len(errStr) > 0 {
		errs = append(errs, errStr)
	}
	return errs
}

func processFormField(r *http.Request, field string) (string, string) {
	fieldData := r.PostFormValue(field)
	if len(fieldData) == 0 {
		return "", "Missing '" + field + "' parameter, cannot continue"
	}
	return fieldData, ""
}
