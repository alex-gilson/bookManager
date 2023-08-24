package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	Genre         string `json:"genre"`
}

func main() {

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "books",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to mysql server.")

	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookByID)
	router.POST("/books", postBooks)
	router.DELETE("/books/:id", deleteBookByID)
	router.Run(":8080")
}

func postBooks(c *gin.Context) {

	var newBook Book
	var err error

	if err = c.BindJSON(&newBook); err != nil {
		return
	}

	if newBook.Author == "" || newBook.Title == "" || newBook.Genre == "" {
		c.IndentedJSON(http.StatusBadRequest, "Missing required fields")
		return
	}

	_, err = time.Parse("2006-01-02", newBook.PublishedDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid published_date format")
		return
	}

	_, err = db.Exec("INSERT INTO book (title, author, published_date, genre) VALUES (?, ?, ?, ?)", newBook.Title, newBook.Author, newBook.PublishedDate, newBook.Genre)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}

	c.IndentedJSON(http.StatusCreated, newBook)

}

func getBookByID(c *gin.Context) {

	id := c.Param("id")
	var book Book

	row := db.QueryRow("SELECT * FROM book WHERE id = ?", id)
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Genre); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func getBooks(c *gin.Context) {
	var books []Book

	query := "SELECT * FROM book WHERE 1=1"
	params := c.Request.URL.Query()

	if id := params.Get("id"); id != "" {
		query += " AND id = " + id
	}

	if title := params.Get("title"); title != "" {
		query += " AND title = \"" + title + "\""
	}

	if author := params.Get("author"); author != "" {
		query += " AND author = \"" + author + "\""
	}

	if genre := params.Get("genre"); genre != "" {
		query += " AND genre = \"" + genre + "\""
	}

	if fromDate := params.Get("published_from"); fromDate != "" {
		_, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid published_from format")
			return
		}
		query += " AND published_date >= \"" + fromDate + "\""
	}

	if toDate := params.Get("published_to"); toDate != "" {
		query += " AND published_date <= \"" + toDate + "\""
		_, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid published_to format")
			return
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Genre); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}

	if len(books) == 0 {
		c.IndentedJSON(http.StatusOK, []Book{})
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

func deleteBookByID(c *gin.Context) {

	var rowsAffected int64

	id := c.Param("id")

	res, err := db.Exec("DELETE FROM book WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
	} else if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
	}
}
