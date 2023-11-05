package utils

import (
	"encoding/json"
	"fmt"
	"github.com/wexinc/ps-tag-onboarding-go/internal/log"
	"net/http"
)

// ParseJson gets json for request and fills the target model
func ParseJson(r *http.Request, target interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&target); err != nil {
		return err
	}

	log.Info.Printf("Target value decoded: %v", target)

	return nil
}

// ResponseJson makes the response with payload as json format
func ResponseJson(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Error.Println(err)
		ResponseError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(response)
}

// ResponseCustomError makes the error response with default message in json format
func ResponseError(w http.ResponseWriter, code int) {
	var message string
	switch code {
	case http.StatusBadRequest:
		message = "The request had invalid inputs or otherwise cannot be served."
	case http.StatusUnauthorized:
		message = "Authorization information is missing or invalid."
	case http.StatusNotFound:
		message = "Unable to find requested record."
	case http.StatusRequestTimeout:
		message = "Request took too long to process."
	case http.StatusRequestedRangeNotSatisfiable:
		message = "No resource available, unable to fulfill the request."
	case http.StatusTooManyRequests:
		message = "Request rate too high, requests from this this user are throttled."
	case http.StatusInternalServerError:
		message = "An error was encountered."
	case http.StatusServiceUnavailable:
		message = "The service is unavailable, please try again later."
	case http.StatusGatewayTimeout:
		message = "The service timed out waiting for an upstream response. Try again later."
	}

	ResponseCustomError(w, code, message)
}

// ResponseCustomError makes the error response with given message in json format
func ResponseCustomError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"code":%v,"message":"%v"}`, code, message)))
}

// ResponseCustomError makes the error response with given message in json format
func ResponseMessageErr(w http.ResponseWriter, msgErr MessageErr) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(msgErr.Status())
	w.Write([]byte(fmt.Sprintf(`{"status":%v,"message":"%v","error":"%v"}`, msgErr.Status(), msgErr.Message(), msgErr.Error())))
}
