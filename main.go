package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// book represents a book with its essential details.
// The struct uses JSON tags to specify the field names when the struct is
// marshaled or unmarshaled from JSON. These tags ensure that the JSON
// representation of the struct uses the specified names, making it easier
// to work with external systems or APIs that rely on JSON data.
type book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// books is a slice of book structs
var books = []book{
	{ID: 1, Title: "The Go Programming Language", Author: "Brian Kernighan", Quantity: 2},
	{ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Quantity: 5},
	{ID: 3, Title: "Head First Go", Author: "Jay McGavren", Quantity: 6},
}

//   - c: A pointer to the Gin context, which contains information about the
//     HTTP request and is used to construct the response.
//
// The function retrieves the `books` slice and sends it as a JSON response
// to the client.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// createBooks handles the HTTP request to create a new book.
// It expects a JSON payload representing a book, which is bound to a `book` struct.
//
// The function performs the following steps:
// 1. Attempts to bind the JSON request body to the `newBook` variable.
// 2. If binding fails (due to invalid JSON), it responds with a 400 Bad Request status and an error message.
// 3. If binding is successful, the new book is appended to the `books` slice.
// 4. Responds with a 201 Created status and the newly created book in the response body.

func createBooks(c *gin.Context) {

	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// bookById retrieves a book by its ID from the URL parameter and returns it as a JSON response.
// If the ID is invalid or if the book is not found, it responds with an appropriate HTTP status code and error message.
func bookById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

// getBookById searches for a book in the books slice by its ID and returns the book if found.
// If the book is not found, it returns an error indicating that the book was not found.
//
// @param id int - The ID of the book to search for.
// @return (*book, error) - A pointer to the book if found, or nil if not found, along with an error indicating the result of the search.
func getBookById(id int) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

// checkoutBook handles the checkout process for a book.
// It expects an "id" query parameter in the request URL, which represents the ID of the book to be checked out.
//
// The function performs the following steps:
// 1. Retrieves the "id" query parameter from the request.
// 2. If the "id" parameter is missing, it responds with a 400 Bad Request status and a message indicating the missing parameter.
// 3. Converts the "id" parameter from a string to an integer. If the conversion fails, it responds with a 400 Bad Request status and an error message.
// 4. Fetches the book details using the provided ID. If the book is not found, it responds with a 404 Not Found status and a message indicating the book was not found.
// 5. Checks if the book's quantity is greater than zero. If the book is out of stock, it responds with a 400 Bad Request status and a message indicating the book is not available.
// 6. Decreases the book's quantity by one to reflect the checkout action.
//

func checkoutBook(c *gin.Context) {
	idStr, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query parameter"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

// Creates a new Gin router instance with default middleware.
// Registers the `getBooks` handler function to the "/books" route. This
// route will handle GET requests to retrieve a list of books.
func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBooks)
	router.GET("/books/:id", bookById)
	router.GET("/checkout", checkoutBook)
	router.Run("localhost:8080")

}
