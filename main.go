package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

//For public usage outside this module, we can only access structs,
// if the attributes are in capital

// Reason why we mention json: "id" is cause in an occurence of serialization,
// json will look for it via the simple letter terms assigned next to the attributes
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) { //pointer to a context -> all information about the request
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil { //if there is an error, it will simply return
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// This will return the pointer for the bookId
func getBookById(id string) (*book, error) {
	for i, b := range books { //This works as a foreach(books in book) where i is for index and b is for book
		if b.ID == id {
			return &books[i], nil
		}
	}

	//This will get triggered if the for loop will be run out without a return
	return nil, errors.New("Book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Missing id query parameter"})
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity = book.Quantity - 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	//http://localhost:8080/checkout?id=3

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Wrong book inputted. Try again"})
		return
	}

	book.Quantity = book.Quantity + 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/book", getBooks) //the route
	router.POST("/book", createBook)
	router.GET("/book/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080") //this will be port
}
