package main

import (
	// "encoding/json"
	"log"
	"fmt"
	"net/http"
	"context"
	"sync"
	"time"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client
//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error
//Used to execute client creation procedure only once.
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
	message string
}

func getClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI("mongodb+srv://dbuser:YlxiFoOkwEWnwgYt@cluster0.rhw7i.mongodb.net/")
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
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

func getBobaCount() (*Username) {
	user := Username{}
	return &user
}

func getUsersInServer(server string) []string {
	client, err := getClient()
	if err != nil {
		return nil
	}

	cursor, err := client.Database("boba_db").Collection(server).Find(context.Background(), bson.M{})
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
		if el == server {
			return true
		}
	}
	return false
}

func getBoba(w http.ResponseWriter, r *http.Request) {
	server := mux.Vars(r)["server"]
	log.Printf("GET /boba/%s", server)
	fmt.Fprint(w, "server")
	fmt.Println(server)

	if !doesServerExist(server) {
		fmt.Println("server doesn't exist error")
		// return json for server doesn't exist
	}

	userList := getUsersInServer(server)
	fmt.Println(userList)


	client, err := getClient()
	if err != nil {
		// return json connection error
	}

	filter := bson.M{"user": bson.D{primitive.E{Key: "$in", Value: userList}}}
	cursor, err := client.Database("boba_db").Collection("boba_count").Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var users []UserData

	if err = cursor.All(context.Background(), &users); err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)

	// fmt.Println(&client)
	// else return json of users and 

	// w.Header().Set("Content-Type", "application/json")
	// TODO: return sorted json by values
}

func main() {
	godotenv.Load()

	router := mux.NewRouter()
	router.HandleFunc("/boba/{server}", getBoba).Methods("GET")
	router.Use(mux.CORSMethodMiddleware(router))

	// exist := doesServerExist("Rieber Heauxs")
	// fmt.Println(exist)
	// fmt.Println("hi")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Listening on port 8000")
	_, _ = getClient()
	log.Fatal(srv.ListenAndServe())
}
