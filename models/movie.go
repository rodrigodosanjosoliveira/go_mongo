package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Movie struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Movie  string             `json:"movie"`
	Actors []string           `json:"actors"`
}

func InsertMovie(movie Movie) error {
	inserted, err := databaseCollection().InsertOne(context.TODO(), movie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a record with id: ", inserted.InsertedID)

	return err
}

func InsertMany(movies []Movie) error {
	newMovies := make(any, len(movies))
	for i, movie := range movies {
		newMovies[i] = movie
	}

	result, err := databaseCollection().InsertMany(context.TODO(), newMovies)
	if err != nil {
		panic(err)
	}

	log.Println(result)

	return err
}

func UpdateMovie(movieID string, movie Movie) error {
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"movie": movie.Movie, "actors": movie.Actors}}

	result, err := databaseCollection().UpdateByID(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Println("New record: ", result)

	return nil
}

func DeleteMovie(movieID string) error {
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	result, err := databaseCollection().DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	fmt.Println("Record deleted: ", result)

	return err
}

func Find(movieName string) Movie {
	var result Movie

	filter := bson.D{{Key: "movie", Value: movieName}}

	err := databaseCollection().FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func FindAll(movieName string) []Movie {
	var results []Movie

	cursor, err := databaseCollection().Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	err = cursor.All(context.TODO(), results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func DeleteAll() error {
	delResult, err := databaseCollection().DeleteMany(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		return err
	}

	fmt.Println("Records deleted:", delResult.DeletedCount)

	return err
}

func databaseCollection() *mongo.Collection {
	return mongoClient.Database(db).Collection(collName)
}
