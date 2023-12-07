package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Price     int       `json:"price,omitempty" bson:"price,omitempty"`
	CreatedOn time.Time `json:"createdon,omitempty" bson:"createdon,omitempty"`
}

type BookStore struct {
	C      *mongo.Collection
	Client *mongo.Client
	Table  string
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func NewDB(table string) (*BookStore, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	collection := client.Database("store").Collection(table)
	books := BookStore{
		Client: client,
		C:      collection,
		Table:  table,
	}
	return &books, nil
}

func (books BookStore) GetDBName() string {
	return books.Table
}

func (books BookStore) AddBook(title string, price int) error {

	b := Book{
		Title:     title,
		Price:     price,
		CreatedOn: time.Now(),
	}
	doc, err := toDoc(b)
	if err != nil {
		return err
	}
	res, err := books.C.InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("record book <%s> inserted with id = %s\n", b.Title, res.InsertedID)
	return err
}

func printBooks(books []Book) {
	for _, book := range books {
		fmt.Printf("Title: %s, Price: %d\n", book.Title, book.Price)
	}
}
func (books BookStore) DeleteBook(title string) {
	//filter := bson.D{{"title", title}}
	_, err := books.C.DeleteOne(context.Background(), bson.M{"title": title})
	if err != nil {
		fmt.Printf("Error deleting book %s\n", title)
		return
	}
	fmt.Printf("Book %s deleted\n", title)
}

func (books BookStore) GetAll() []Book {
	var b []Book
	cur, err := books.C.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		result := Book{}
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		b = append(b, result)
	}
	printBooks(b)
	return b
}

func (books BookStore) Disconnect() {
	fmt.Println("Closing Bookstore")
	books.Client.Disconnect(context.Background())
}
