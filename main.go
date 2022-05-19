package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

//database connection

// get request to homepage
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func main() {

	username := "sarahwatremet"
	password := "dpxy0cZJYxvk8RUm"
	cluster := "cluster-de-fifou.bxtvr.mongodb.net"

	uri := "mongodb+srv://" + url.QueryEscape(username) + ":" +
		url.QueryEscape(password) + "@" + cluster
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database("hawaii-surf-spots").Collection("surf-spots")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/spot", createSpot).Methods("POST")
	router.HandleFunc("/spots", getAllSpots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOneSpot).Methods("GET")
	// router.HandleFunc("/spots/{id}", updateSpot).Methods("PATCH")
	// router.HandleFunc("/spots/{id}", deleteSpot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// create the strcture of the database in JSON

type Spots struct {
	Records []Record `json:"records,omitempty"`
}

type Record struct {
	ID          string  `json:"id,omitempty"`
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

var spots Spots

// post request createSpot
func createSpot(w http.ResponseWriter, r *http.Request) {

	var newSpot Record
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the spot only in order to update")
	}

	json.Unmarshal(reqBody, &newSpot)
	spots.Records = append(spots.Records, newSpot)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newSpot)
}

// get request for one spot
func getOneSpot(w http.ResponseWriter, r *http.Request) {
	// opened spot.json
	jsonFile, err := os.Open("spot.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened spot.json")

	defer jsonFile.Close()

	// decrypted json file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// initialized with variable according to the struct
	var spots Spots

	// get json informations to the variable
	json.Unmarshal(byteValue, &spots)

	spotID := mux.Vars(r)["id"]

	for _, singleSpot := range spots.Records {
		if singleSpot.ID == spotID {
			json.NewEncoder(w).Encode(singleSpot)
		}
	}
}

// get request for all spots
func getAllSpots(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("spot.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened spot.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &spots)
	json.NewEncoder(w).Encode(spots)
}

// // update request for an spot
// func updateSpot(w http.ResponseWriter, r *http.Request) {
// 	spotID := mux.Vars(r)["id"]
// 	var updatedSpot spot

// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the spot title and description only in order to update")
// 	}
// 	json.Unmarshal(reqBody, &updatedSpot)

// 	for i, singleSpot := range spots {
// 		if singleSpot.ID == eventID {
// 			singleSpot.Title = updatedEvent.Title
// 			singleSpot.Description = updatedEvent.Description
// 			spots = append(spots[:i], singleSpot)
// 			json.NewEncoder(w).Encode(singleSpot)
// 		}
// 	}
// }

// // delete request to remove an spot
// func deleteSpot(w http.ResponseWriter, r *http.Request) {
// 	spotID := mux.Vars(r)["id"]

// 	for i, singleSpot := range spots {
// 		if singleSpot.ID == spotID {
// 			spots = append(spots[:i], spots[i+1:]...)
// 			fmt.Fprintf(w, "The spot with ID %v has been deleted successfully", spotID)
// 		}
// 	}
// }
