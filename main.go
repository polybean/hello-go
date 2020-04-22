package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var serviceName = "hello-go"
var coll *mgo.Collection
var histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Subsystem: "http_server",
	Name:      "resp_time",
	Help:      "Request response time",
}, []string{
	"service",
	"code",
	"method",
	"path",
})

type greeting struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Message string        `bson:"message" json:"message"`
}

func init() {
	prometheus.MustRegister(histogram)
}

func hello(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	defer func() { recordMetrics(start, request, http.StatusOK) }()

	log.Printf("%s request to %s\n", request.Method, request.RequestURI)
	delay := request.URL.Query().Get("delay")

	if len(delay) > 0 {
		delayNum, _ := strconv.Atoi(delay)
		time.Sleep(time.Duration(delayNum) * time.Millisecond)
	}
	io.WriteString(writer, "Hello, Go v3.0!\n")
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

func version(writer http.ResponseWriter, request *http.Request) {
	log.Printf("%s request to %s\n", request.Method, request.RequestURI)
	release := request.Header.Get("release")
	if release == "" {
		release = "unknown"
	}
	msg := fmt.Sprintf("Version: %s; Release: %s\n", os.Getenv("VERSION"), release)
	io.WriteString(writer, msg)
}

func randomError(writer http.ResponseWriter, request *http.Request) {
	code := http.StatusOK
	start := time.Now()
	defer func() { recordMetrics(start, request, code) }()

	log.Printf("%s request to %s\n", request.Method, request.RequestURI)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)
	msg := "Everything is still OK"
	version := os.Getenv("VERSION")
	if len(version) > 0 {
		msg = fmt.Sprintf("%s with version %s", msg, version)
	}
	msg = fmt.Sprintf("%s\n", msg)
	if n == 0 {
		code = http.StatusInternalServerError
		msg = "ERROR: Something, somewhere, went wrong!\n"
		log.Printf(msg)
	}
	writer.WriteHeader(code)
	io.WriteString(writer, msg)
}

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/greetings", greetings)
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/random-error", randomError)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", version)

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

// Add function is just for unit test demonstration
func Add(x, y int) int {
	return x + y
}

func main() {
	log.Printf("Starting the application\n")
	if len(os.Getenv("SERVICE_NAME")) > 0 {
		serviceName = os.Getenv("SERVICE_NAME")
	}
	connect2Mongo()
	startServer()
}

func recordMetrics(start time.Time, req *http.Request, code int) {
	duration := time.Since(start)
	histogram.With(
		prometheus.Labels{
			"service": serviceName,
			"code":    fmt.Sprintf("%d", code),
			"method":  req.Method,
			"path":    req.URL.Path,
		},
	).Observe(duration.Seconds())
}
