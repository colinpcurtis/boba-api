package main

import (
	"os"
	"log"
	"fmt"
	"sort"
	"sync"
	"time"
	"context"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once
var ctx = context.TODO()


type Username struct {
	// must be capitalized
    ID primitive.ObjectID `bson:"_id"`
	User string `bson:"user"`
}

type UserData struct {
	ID primitive.ObjectID `bson:"_id"`
	User string `bson:"user"`
	BobaCount int `bson:"boba_count"`
}

type Error struct {
	Message string
}

type JsonPayload struct {
	User string
	BobaCount int
}

func getClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
		log.Println("Connected to MongoDB")
	})
	return clientInstance, clientInstanceError
}

func prepareJson(users []UserData) ([]JsonPayload) {
	var json []JsonPayload
	for _, user := range(users) {
		data := JsonPayload{user.User, user.BobaCount}
		json = append(json, data)
	}
	return json
}

func getUsersInServer(server string) []string {
	cursor, err := clientInstance.Database("boba_db").Collection(server).Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	
	var results []Username
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	var users []string
	for _, userObject := range results {
		users = append(users, userObject.User)
	}
	return users
}

func doesServerExist(server string) bool {
	client, err := getClient()
	if err != nil {
		fmt.Println("err")
		return false
	}
	serverList, err := client.Database("boba_db").ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return false
	}
	for _, el := range serverList {
		if el != "boba_count" && el == server {
			return true
		}
	}
	return false
}

func getBoba(w http.ResponseWriter, r *http.Request) {
	server := mux.Vars(r)["server"]
	log.Printf("GET /boba/%s", server)
	w.Header().Set("Content-Type", "application/json")

	if !doesServerExist(server) {
		message := "Error: " + server + " does not exist, please try again"
		error := Error{message}
		json.NewEncoder(w).Encode(error)
		return
	}
	userList := getUsersInServer(server)

	filter := bson.M{"user": bson.D{primitive.E{Key: "$in", Value: userList}}}
	cursor, err := clientInstance.Database("boba_db").Collection("boba_count").Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var users []UserData

	if err = cursor.All(context.Background(), &users); err != nil {
		log.Fatal(err)
	}
	
	payload := prepareJson(users)
	sort.Slice(payload, func(i, j int) bool {
		return payload[i].BobaCount > payload[j].BobaCount
	})
	json.NewEncoder(w).Encode(payload)
	return
}

func main() {
	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/boba/{server}", getBoba).Methods("GET")
	router.Use(mux.CORSMethodMiddleware(router))

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Listening on port 8000")
	_, _ = getClient()
	log.Fatal(srv.ListenAndServe())
}
