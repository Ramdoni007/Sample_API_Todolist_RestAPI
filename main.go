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

type event struct {
	ID        string `json:"ID"`
	Todo      string `json:"Todo"`
	Completed bool   `json:"Completed"`
}

type allEvents []event

var events = allEvents{
	{
		ID:        "1",
		Todo:      "Mempunyai Pacar Jenny BlackPink",
		Completed: false,
	},
	{
		ID:        "2",
		Todo:      "Mempunyai Pacar Jenny BlackPink",
		Completed: false,
	},
	{
		ID:        "3",
		Todo:      "Mempunyai Pacar Jenny BlackPink",
		Completed: true,
	},
	{
		ID:        "4",
		Todo:      "Mempunyai Pacar Jenny BlackPink",
		Completed: true,
	},
	{
		ID:        "5",
		Todo:      "Mempunyai Pacar Jenny BlackPink",
		Completed: false,
	},
	{
		ID:        "6",
		Todo:      "Mempunyai Pacar Jenny BlackPink",
		Completed: false,
	},
}

func homePoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home")
}

func createEvent(w http.ResponseWriter, r *http.Request) {

	var newEvent event

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Jika Kita ingin menambahkan Todo silahkan masukan Todo yang ingin anda buat ")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)

}

func getOnEvent(w http.ResponseWriter, r *http.Request) {

	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvent(w http.ResponseWriter, r *http.Request) {
	completed := r.URL.Query().Get("completed")
	offSet := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	var newEvents allEvents

	if completed == "true" {
		for _, singleEvent := range events {
			if singleEvent.Completed {
				newEvents = append(newEvents[:], singleEvent)

			}
		}
	} else if completed == "false" {
		for _, singleEvent := range events {
			if !singleEvent.Completed {
				newEvents = append(newEvents[:], singleEvent)

			}
		}
	} else {
		for _, singleEvent := range events {
			newEvents = append(newEvents[:], singleEvent)

		}
	}

	var limiter, _ = strconv.Atoi(limit)
	var offseter, _ = strconv.Atoi(offSet)

	if limiter < 0 || limiter > len(newEvents) {
		limiter = 0
	}

	if offseter == 0 && limiter > offseter {
		offseter = limiter
	}

	if offseter == 0 || offseter > len(newEvents) {

	}

	if limiter == 0 {
		limiter = offseter
	}

	json.NewEncoder(w).Encode(newEvents[offseter-limiter : offseter])
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Jika Kita ingin Menambahkan Todo silahkan masukan todo yang ingin anda buat")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Todo = updatedEvent.Todo
			singleEvent.Completed = updatedEvent.Completed
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}

}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "Data Dengan ID %v telah berhasil di hapus ", eventID)
		}
	}
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePoint)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event", getAllEvent).Methods("GET")
	router.HandleFunc("/event/{id}", getOnEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/event/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

}
