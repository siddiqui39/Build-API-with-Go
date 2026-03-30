
package main 

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)
//Define a struct with public fields; uppercase
type book struct{
	ID		string	`json:"id"`		
	Title 	string 	`json:"title"`
	Author 	string	`json:"author"`
	Quantity  int 	`json:"quantity"`
}

//List of books
var books = []book{
    {ID: "101", Title: "A Tale of Two Cities", Author: "Charles Dickens", Quantity: 2},
    {ID: "102", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 33},
    {ID: "103", Title: "12 Rules for Life", Author: "Jordan B Peterson", Quantity: 12},
}

//Returns all books as JSON
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)

}

//Returns a single book by ID
func bookById(c *gin.Context){
	id := c.Param("id")				//Get book ID from URL
	book, err := getBookById(id)	//Find book

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// Send the book as JSON
	c.IndentedJSON(http.StatusOK, book)

}

//Checks out a book
func checkoutBook(c *gin.Context){
	id, ok := c.GetQuery("id")		//Get book ID from query parameter

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)	//Find book
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {			//check if available
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1				//Reduce quantity
	c.IndentedJSON(http.StatusOK, book)
}

//Returns a book
func  returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")		//Get ID from query parameter

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)	//Find book
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1				//Increase quantity
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error){
	for i, b := range books {
		if b.ID == id{
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

//Create a new book
func createBook(c*gin.Context){
	var newBook book 

	//Bind JSON from request body to newBook
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	books= append(books, newBook)		//Add book to list
	c.IndentedJSON(http.StatusCreated, newBook)		//Return created book
}

//Main function to setup routes and start server
func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}


