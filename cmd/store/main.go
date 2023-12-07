package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vlakshmi70/store/internal/db"
)

func init() {
	fmt.Println("Starting Bookstore app")
}

func main() {
	mydb, err := db.NewDB("books")
	if err != nil {
		log.Fatal(err)
	}
	defer mydb.Disconnect()
	var choice, price int
	var title string
	for {
		fmt.Printf("Enter 1 to add-books, 2 to del-books, 3 to display, 4 to exit: ")
		fmt.Scanf("%d", &choice)
		fmt.Println()
		switch choice {
		case 1:
			fmt.Printf("Enter title: ")
			fmt.Scanf("%s", &title)
			fmt.Println()
			fmt.Printf("Enter price: ")
			fmt.Scanf("%d", &price)
			fmt.Println()
			mydb.AddBook(title, price)
			break
		case 2:
			fmt.Printf("Enter title: ")
			fmt.Scanf("%s", &title)
			fmt.Println()
			mydb.DeleteBook(title)
		case 3:
			mydb.GetAll()
			break
		case 4:
			fmt.Println("Exiting bookstore app")
			os.Exit(0)
			break
		default:
			fmt.Println("Wrong choice, try again")
		}
	}
}
