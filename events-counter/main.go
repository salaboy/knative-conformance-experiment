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

var receivedEvents map[string]*event.Event
var failedEvents map[string]*event.Event

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/events", EventReportHandler).Methods("GET")
	r.HandleFunc("/events", EventReceiverHandler).Methods("POST")
	r.HandleFunc("/events", EventDeleteReceiverHandler).Methods("DELETE")
	r.HandleFunc("/events/data-plane/delivery-retry", EventDeliveryRetryReceiverHandler).Methods("POST")
	r.HandleFunc("/events/data-plane/delivery-retry/report", EventDeliveryRetryReportReceiverHandler).Methods("GET")

	receivedEvents = make(map[string]*event.Event)
	failedEvents = make(map[string]*event.Event)
	log.Printf("Events Counter 8080!")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func EventReportHandler(writer http.ResponseWriter, request *http.Request) {
	respondWithJSON(writer, http.StatusOK, &receivedEvents)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func EventDeleteReceiverHandler(writer http.ResponseWriter, request *http.Request) {
	receivedEvents = make(map[string]*event.Event)
	failedEvents = make(map[string]*event.Event)
	fmt.Printf("Reseting Event Store ...")
	respondWithJSON(writer, http.StatusOK, nil)
}

func EventReceiverHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := context.Background()
	message := cehttp.NewMessageFromHttpRequest(request)
	event, _ := binding.ToEvent(ctx, message)
	receivedEvents[event.ID()] = event
	fmt.Printf("Got an Event: %s", event)
	respondWithJSON(writer, http.StatusOK, &event)
}

var failForRedeliveryFlag = true

func EventDeliveryRetryReceiverHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	message := cehttp.NewMessageFromHttpRequest(request)
	event, _ := binding.ToEvent(ctx, message)
	fmt.Printf("Got an Event: %s\n", event)
	fmt.Printf("Should I fail? %s\n", failForRedeliveryFlag)
	if failForRedeliveryFlag {
		fmt.Printf("But I am returning: %s\n", http.StatusBadRequest)
		failedEvents[event.ID()] = event
		respondWithJSON(writer, http.StatusBadRequest, &event)
		failForRedeliveryFlag = false
	} else {
		receivedEvents[event.ID()] = event
		fmt.Printf("I am returning: %s\n", http.StatusOK)
		respondWithJSON(writer, http.StatusOK, &event)
	}

}
type EventDeliveryRetryReport struct{
	receivedEvents map[string]*event.Event
	failedEvents map[string]*event.Event
}
func EventDeliveryRetryReportReceiverHandler(writer http.ResponseWriter, request *http.Request) {
	var report = EventDeliveryRetryReport{
		receivedEvents: receivedEvents,
		failedEvents: failedEvents,
	}
	respondWithJSON(writer, http.StatusOK, &report)
}