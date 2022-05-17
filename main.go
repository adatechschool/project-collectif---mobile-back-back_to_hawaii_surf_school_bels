package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// get request to homepage
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func main() {
	jsonFile, err := os.Open("spot.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened spot.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// create the strcture of the database in JSON

type Event struct {
	Records []Record `json:"records,omitempty"`
}

type Record struct {
	ID          *string `json:"id,omitempty"`         
	CreatedTime *string `json:"createdTime,omitempty"`
	Fields      *Fields `json:"fields,omitempty"`     
}

type Fields struct {
	SurfBreak               []string `json:"Surf Break,omitempty"`               
	DifficultyLevel         *int64   `json:"Difficulty Level,omitempty"`         
	Destination             *string  `json:"Destination,omitempty"`              
	Geocode                 *string  `json:"Geocode,omitempty"`                  
	Influencers             []string `json:"Influencers,omitempty"`              
	MagicSeaweedLink        *string  `json:"Magic Seaweed Link,omitempty"`       
	Photos                  []Photo  `json:"Photos,omitempty"`                   
	PeakSurfSeasonBegins    *string  `json:"Peak Surf Season Begins,omitempty"`  
	DestinationStateCountry *string  `json:"Destination State/Country,omitempty"`
	PeakSurfSeasonEnds      *string  `json:"Peak Surf Season Ends,omitempty"`    
	Address                 *string  `json:"Address,omitempty"`                  
}

type Photo struct {
	ID         *string     `json:"id,omitempty"`        
	Width      *int64      `json:"width,omitempty"`     
	Height     *int64      `json:"height,omitempty"`    
	URL        *string     `json:"url,omitempty"`       
	Filename   *string     `json:"filename,omitempty"`  
	Size       *int64      `json:"size,omitempty"`      
	Type       *string     `json:"type,omitempty"`      
	Thumbnails *Thumbnails `json:"thumbnails,omitempty"`
}

type Thumbnails struct {
	Small *Full `json:"small,omitempty"`
	Large *Full `json:"large,omitempty"`
	Full  *Full `json:"full,omitempty"` 
}

type Full struct {
	URL    *string `json:"url,omitempty"`   
	Width  *int64  `json:"width,omitempty"` 
	Height *int64  `json:"height,omitempty"`
}

type allEvents []event

// database according to the structure
var events = allEvents{
	{
		id: "1",
		SurfBreak: ["Reef Break"],
		DifficultyLevel: "4",
		Destination: "Hawaii",
		Coordinates: {
			Latitude: "0",
			Longitude: "1",
		},
		MagicSeaweed: "https://fr.magicseaweed.com/",
		Photos: "https://images.unsplash.com/photo-1502680390469-be75c86b636f?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		SeasonStart: "27 juin 2022",
		SeasonEnd: "23 octobre 2022",
		Address: "île de Hawaii",
	},
}



// post request createEvent
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

// get request for one event
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

// get request for all events
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

// update request for an event
func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

// delete request to remove an event
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}
