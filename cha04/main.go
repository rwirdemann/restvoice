package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Contact struct {
}

var contacts = make(map[int]string)
var contactMap = make(map[int]Contact)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/contacts/{id:[0-9]+}", getContact).Methods("GET")
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
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(b)
}

func addContactHandler(writer http.ResponseWriter, request *http.Request) {
	if contact, err := ioutil.ReadAll(request.Body); err == nil {
		if len(contact) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		id := nextId()
		contacts[id] = string(contact)
		url := request.URL.String()
		writer.Header().Set("Location", fmt.Sprintf("%s/%d", url, id))
		writer.WriteHeader(http.StatusCreated)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func updateContactHandler(writer http.ResponseWriter, request *http.Request) {
	if contact, err := ioutil.ReadAll(request.Body); err == nil {
		if len(contact) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		v := mux.Vars(request)
		id, _ := strconv.Atoi(v["id"])
		contacts[id] = string(contact)
		writer.WriteHeader(http.StatusNoContent)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func getContact(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	writer.Write([]byte(contacts[id]))
}

func deleteContactHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	if _, ok := contacts[id]; !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	delete(contacts, id)
	writer.WriteHeader(http.StatusNoContent)
}

func nextId() int {
	id := 1
	for k, _ := range contacts {
		if k > id {
			id = k + 1
		}
	}
	return id
}
