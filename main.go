package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

type book struct {
	Id        string `json:"id"`   
	Title     string `json:"title"`
	Author    string `json:"author"`
	Quantity  int    `json:"quantity"`
}

var books = []book{
	{Id: "1", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 3},
	{Id: "2", Title: "To Kill a Mockingbird", Author: "Harper Lee", Quantity: 2},
	{Id: "3", Title: "1984", Author: "George Orwell", Quantity: 5},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getName(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}

func createNewBook(c *gin.Context) {
	var newBook book
	 
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*book, error){

	for index, b := range books {
		if b.Id == id {
			return &books[index], nil
		}
	}

	return nil, errors.New("book not found")
}

func BookIdRouter(c *gin.Context) {
	id := c.Param("id")

	book, err := getBookById(id);

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func CheckOutBook(c *gin.Context) {

	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	book, err := getBookById(id);

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if book.Quantity == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "book is out of stock"})
		return
	}

	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)
}

func deleteBookById(id string) error {
	
	for index, b := range books {
		if b.Id == id {
			books = append(books[:index], books[index+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}

func RouterDelete(c *gin.Context) {
	id := c.Param("id")

	err := deleteBookById(id);

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted"})

}

func main() {
	router := gin.Default();
	router.GET("/books", getBooks)
	router.GET("/", getName)
	router.POST("/createNewBook", createNewBook)
	router.GET("/books/:id", BookIdRouter) 
	router.PATCH("/checkout", CheckOutBook)
	router.DELETE("/delBooks/:id", RouterDelete)
	router.Run("localhost:8080")
}
