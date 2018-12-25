package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Contact struct {
	Firstname string
	Lastname  string
}

var contactMap = make(map[int]Contact)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/contacts", getContactsHandler).Methods("GET")
	r.HandleFunc("/contacts/{id:[0-9]+}", getContactHandler).Methods("GET")
	r.HandleFunc("/contacts", addContactHandler).Methods("POST")
	r.HandleFunc("/contacts/{id[0-9]+}", updateContactHandler).Methods("PUT")
	r.HandleFunc("/contacts/{id[0-9]+}", deleteContactHandler).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}

func getContactsHandler(writer http.ResponseWriter, request *http.Request) {
	var s []Contact
	for _, c := range contactMap {
		s = append(s, c)
	}
	b, err := json.Marshal(s)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(b)
}

func getContactHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	if _, ok := contactMap[id]; !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(contactMap[id])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(b)
}

func addContactHandler(writer http.ResponseWriter, request *http.Request) {
	if b, err := ioutil.ReadAll(request.Body); err == nil {
		if len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		id := nextId()
		var contact Contact
		json.Unmarshal(b, &contact)
		contactMap[id] = contact
		url := request.URL.String()
		writer.Header().Set("Location", fmt.Sprintf("%s/%d", url, id))
		writer.WriteHeader(http.StatusCreated)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func updateContactHandler(writer http.ResponseWriter, request *http.Request) {
	if b, err := ioutil.ReadAll(request.Body); err == nil {
		if len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		v := mux.Vars(request)
		id, _ := strconv.Atoi(v["id"])
		var contact Contact
		json.Unmarshal(b, &contact)
		contactMap[id] = contact
		writer.WriteHeader(http.StatusNoContent)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func deleteContactHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	if _, ok := contactMap[id]; !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	delete(contactMap, id)
	writer.WriteHeader(http.StatusNoContent)
}

func nextId() int {
	id := 1
	for k, _ := range contactMap {
		if k >= id {
			id = k + 1
		}
	}
	return id
}
