package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var events map[string]*event.Event

var failFlag = true

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/events", EventReportHandler).Methods("GET")
	r.HandleFunc("/events", EventReceiverHandler).Methods("POST")
	r.HandleFunc("/events-fail-once", EventFailOnceReceiverHandler).Methods("POST")

	events = make(map[string]*event.Event)
	log.Printf("Events Counter 8080!")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func EventReportHandler(writer http.ResponseWriter, request *http.Request) {
	respondWithJSON(writer, http.StatusOK, &events)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func EventReceiverHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := context.Background()
	message := cehttp.NewMessageFromHttpRequest(request)
	event, _ := binding.ToEvent(ctx, message)
	events[event.ID()] = event
	fmt.Printf("Got an Event: %s", event)
	respondWithJSON(writer, http.StatusOK, &event)
}

func EventFailOnceReceiverHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := context.Background()
	message := cehttp.NewMessageFromHttpRequest(request)
	event, _ := binding.ToEvent(ctx, message)
	events[event.ID()] = event
	fmt.Printf("Got an Event: %s\n", event)
	fmt.Printf("Should I fail? %s\n", failFlag)
	if failFlag {
		fmt.Printf("But I am returning: %s\n", http.StatusBadRequest)
		respondWithJSON(writer, http.StatusBadRequest, &event)
		failFlag = false
	} else {
		failFlag = true
		fmt.Printf("I am returning: %s\n", http.StatusOK)
		respondWithJSON(writer, http.StatusOK, &event)
	}

}
