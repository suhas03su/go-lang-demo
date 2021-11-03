package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.come/suhas03su/mongoAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017" // THIS should be in ENV
const dbName = "netflix"
const collectionName = "watchlist"

// MOST IMPORTANT
var collection *mongo.Collection

// Connect with MONGO DB
func init() { // RUNS AT THE START and ONLY ONCE
	// CLIENT OPTION
	clientOption := options.Client().ApplyURI(connectionString)

	// CONNECT TO MONGO DB
	client, err := mongo.Connect(context.TODO(), clientOption) // CONTEXT is somewhat similar to CORE DATA context
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MONGO DB Connection SUCCESS")
	collection = client.Database(dbName).Collection(collectionName)

	// IF COLLECTION INSTANCE IS READY
	fmt.Println("COLLECTION REFERENCE IS READY", collection)
}

/*
	MONGODB HELPERS
	INSERT 1 RECORD
*/
func insertOneMovie(movie models.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("INSERTED ONE MOVIE in DB with id -> ", inserted.InsertedID)
}

/*
	MONGODB HELPERS
	UPDATE 1 RECORD
*/
func updateOneMovie(movieId string) {
	objectId, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}

	// FIND ONE AND UPDATE
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"watched": true}}

	updated, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updated Value of the movie -> ", updated.ModifiedCount)
}

/*
	MONGODB HELPERS
	DELETE 1 RECORD
*/
func deleteOneMovie(movieId string) {
	objectId, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}

	// FIND ONE AND DELETE
	filter := bson.M{"_id": objectId}
	deleted, err := collection.DeleteOne(context.Background(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The Deleted movie is -> ", deleted.DeletedCount)
}

/*
	MONGODB HELPERS
	DELETE ALL RECORDS
*/
func deleteAllMovies() int64 {
	// DELETE ALL
	deleted, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The Deleted movie is -> ", deleted.DeletedCount)

	return deleted.DeletedCount
}

/*
	MONGODB HELPERS
	READ ALL RECORDS
*/
func getAllMovies() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}}) // FIND Returns a CURSOR
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M

	// LOOP THROUGH CURSOR
	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	defer cursor.Close(context.Background())

	return movies
}

/*
	ACTUAL CONTROLLERS
	The above one's are DB Helpers
*/
func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateNewMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	defer r.Body.Close()

	var movie models.Netflix
	json.NewDecoder(r.Body).Decode(&movie)

	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	defer r.Body.Close()

	var movie map[string]string
	json.NewDecoder(r.Body).Decode(&movie)

	id := movie["id"]
	updateOneMovie(id)

	json.NewEncoder(w).Encode("The ID sent has been updated!!!!")
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	defer r.Body.Close()

	var movie map[string]string
	json.NewDecoder(r.Body).Decode(&movie)

	id := movie["id"]
	deleteOneMovie(id)

	json.NewEncoder(w).Encode("The Movie with sent ID has been deleted!!!")
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	defer r.Body.Close()

	numberOfMoviesDeleted := deleteAllMovies()

	json.NewEncoder(w).Encode(numberOfMoviesDeleted)
}
