package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Validatable interface {
	Validate() error
}

func splitContentType(ct string) (string, error) {

	arr := strings.Split(ct, "/")
	if len(arr) < 2 {
		msg := "could not split: " + ct
		return "", errors.New(msg)
	}
	l := len(arr)
	return arr[l-1], nil
}

func makeInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err.Error())
		log.Println("makeInt: ", err.Error())
		return 0
	}
	return i
}

func decodeAndValidate(r *http.Request, v Validatable) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(v); err != nil {
		return err
	}

	defer r.Body.Close()

	return v.Validate()
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

/*
	// validations returns error array - we deal with transforming this
	// []errors  to []string to be able to marshal it...
*/
func respondErrorArray(w http.ResponseWriter, code int, errs []error) {
	respondJSON(w, code, makeStringArray(errs))
}

/* utility functions */

func makeStringArray(errors []error) map[string][]string {
	var messages = make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Error()
	}
	msgMap := make(map[string][]string, 0)
	msgMap["errors"] = messages
	return msgMap
}

func CheckJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			log.Println("Wrong Content-Type requested")
			respondError(w, http.StatusNotAcceptable, "Wrong Content-Type header provided")
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
