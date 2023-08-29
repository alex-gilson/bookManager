package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

	// Connect to the database with the name of the database container and it's login details.
	fmt.Println("Connecting to db")
	var err error
	db, err = sql.Open("mysql", "root:mypassword@tcp(db:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// MySQL server isn't fully active yet.
	// Block until connection is accepted.
	for db.Ping() != nil {
		fmt.Println("Attempting connection to db")
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Connected")

	fmt.Println("Creating table")
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS book (
		id              INT AUTO_INCREMENT NOT NULL,
		title           VARCHAR(128) NOT NULL,
		author          VARCHAR(255) NOT NULL,
		published_date  DATE NOT NULL,
		genre           VARCHAR(64) NOT NULL,
		PRIMARY KEY (id)
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/books", getBooks)
		v1.GET("/books/:id", getBookByID)
		v1.POST("/books", postBooks)
		v1.DELETE("/books/:id", deleteBookByID)
	}
	router.Run(":80")
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
