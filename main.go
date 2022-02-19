package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// type Article struct {
// 	Title   string `json:"Title"`
// 	Desc    string `json:"desc"`
// 	Content string `json:"content"`
// }

// type Articles []Article

const uri = "mongodb+srv://ollieadmin:ollie123@cluster0.y0pmv.mongodb.net/blog?retryWrites=true&w=majority"

func handleRequests() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	coll := client.Database("blog").Collection("blogposts")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(mux.CORSMethodMiddleware(myRouter))
	// origins := handlers.AllowedOrigins([]string{"*"})
	myRouter.HandleFunc("/post/{blogtitle}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		vars := mux.Vars(r)
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		blogPosts, err := coll.Find(ctx, bson.D{{"post_title", vars["blogtitle"]}})
		if err != nil {
			log.Fatal(err)
		}
		var episodes []bson.M
		if err = blogPosts.All(ctx, &episodes); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Endpoint hit: all articles")
		json.NewEncoder(w).Encode(episodes)
	}).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		blogPosts, err := coll.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		var episodes []bson.M
		if err = blogPosts.All(ctx, &episodes); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Endpoint hit: all articles")
		json.NewEncoder(w).Encode(episodes)
	}).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {

	handleRequests()

}
