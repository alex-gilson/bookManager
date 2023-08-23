package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type book struct {
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
	router.GET("/books", getBooks2)
	router.GET("/books/:id", getBookByID)
	router.POST("/books", postBooks)
	router.Run("localhost:8080")
}

// func getBooks(c *gin.Context) {

// 	var books []book

// 	rows, err := db.Query("SELECT * FROM book")
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, "")
// 		return
// 	}
// 	defer rows.Close()
// 	// Loop through rows, using Scan to assign column data to struct fields.
// 	for rows.Next() {
// 		var alb book
// 		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Author, &alb.PublishedDate); err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, "")
// 			return
// 		}
// 		books = append(books, alb)
// 	}
// 	if err := rows.Err(); err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, "")
// 		return
// 	}
// 	c.IndentedJSON(http.StatusOK, books)
// }

func postBooks(c *gin.Context) {

	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	_, err := db.Exec("INSERT INTO book (title, author, published_date) VALUES (?, ?, ?)", newBook.Title, newBook.Author, newBook.PublishedDate)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
	}

	c.IndentedJSON(http.StatusCreated, newBook)

}

func getBookByID(c *gin.Context) {

	id := c.Param("id")
	var alb book

	row := db.QueryRow("SELECT * FROM book WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Author, &alb.PublishedDate); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	c.IndentedJSON(http.StatusOK, alb)
}

// booksByAuthor queries for books that have the specified author name.
func getBooksByAuthor(c *gin.Context) {

	name := c.Param("author")

	var books []book

	rows, err := db.Query("SELECT * FROM book WHERE author = ?", name)
	if err != nil {
		// return nil, fmt.Errorf("booksByAuthor %q: %v", name, err)
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb book
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Author, &alb.PublishedDate); err != nil {
			// return nil, fmt.Errorf("booksByAuthor %q: %v", name, err)
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
		books = append(books, alb)
	}
	if err := rows.Err(); err != nil {
		// return nil, fmt.Errorf("booksByAuthor %q: %v", name, err)
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	// return books, nil
	c.IndentedJSON(http.StatusOK, books)
}

func getBooks2(c *gin.Context) {
	var books []book

	// Build the SQL query with optional filters
	query := "SELECT * FROM book WHERE 1=1"
	params := c.Request.URL.Query()

	if title := params.Get("title"); title != "" {
		query += " AND title = ?"
	}

	if author := params.Get("author"); author != "" {
		query += " AND author = ?"
	}

	if genre := params.Get("genre"); genre != "" {
		query += " AND genre = ?"
	}

	if fromDate := params.Get("published_from"); fromDate != "" {
		query += " AND published_date >= ?"
	}

	if toDate := params.Get("published_to"); toDate != "" {
		query += " AND published_date <= ?"
	}

	rows, err := db.Query(query, params.Get("title"), params.Get("author"), params.Get("genre"), params.Get("published_from"), params.Get("published_to"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb book
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Author, &alb.PublishedDate, &alb.Genre); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "")
			return
		}
		books = append(books, alb)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "")
		return
	}

	c.IndentedJSON(http.StatusOK, books)
}
