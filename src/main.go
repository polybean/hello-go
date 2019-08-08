package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var coll *mgo.Collection

type greeting struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Message string        `bson:"message" json:"message"`
}

func index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, Go! v4.0")
}

func findAllGreetings(writer http.ResponseWriter, request *http.Request) {
	var messages []greeting
	coll.Find(bson.M{}).All((&messages))
	body, err := json.Marshal(messages)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		writer.WriteHeader(http.StatusNoContent)
	} else {
		writer.Header().Set(http.CanonicalHeaderKey("content-type"), "application/json")
		writer.Write(body)
	}
}

func insertOneGreeting(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var m greeting
	err := json.NewDecoder(request.Body).Decode(&m)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	m.ID = bson.NewObjectId()
	err = coll.Insert(&m)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// location := fmt.Sprintf("%s/greetings/%s", request.Header.Get("Server"), m.ID.Hex())
	// writer.Header().Set(http.CanonicalHeaderKey("location"), location)
	writer.WriteHeader(http.StatusCreated)
}

func greetings(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		findAllGreetings(writer, request)
		return
	}

	if request.Method == http.MethodPost {
		insertOneGreeting(writer, request)
		return
	}

	writer.WriteHeader(http.StatusNotFound)
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/greetings", greetings)
	http.ListenAndServe(":8080", mux)
	log.Println("Server stated")
}

func connect2Mongo() {
	envKey := "DB_HOST"

	db := os.Getenv(envKey)
	if len(db) == 0 {
		db = "localhost"
	}

	session, err := mgo.Dial(db)
	if err != nil {
		panic(err)
	}

	coll = session.DB("hello").C("greetings")
}

func main() {
	connect2Mongo()
	startServer()
}
