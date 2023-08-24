
# Books API Documentation

## GET /v1/books

**Description:** Retrieve a list of books with optional filtering parameters.  
**Access:** open  
**Operation:** sync  
**Return:** Array containing book objects or error message.

### Response

**Example:**

```json
[
    {
        "id": 1,
        "title": "Book Title",
        "author": "Author Name",
        "published_date": "2023-01-15",
        "genre": "Fiction"
    },
    {
        "id": 2,
        "title": "Another Book",
        "author": "Different Author",
        "published_date": "2022-08-03",
        "genre": "Mystery"
    }
]
```

### Alias Fields

**title:** Name of the book. Case sensitive.  
**author:** Name of the writer of the book. Case sensitive.  
**published_from:** Publishing date from which to include books. Format is YYYY-MM-DD.  
**published_to:** Publishing date until which to include books. Format is YYYY-MM-DD.  
**genre:** Genre of the book. Case sensitive.


## GET /v1/books/{id}

**Description:** Retrieve a book by its ID.  
**Access:** open  
**Operation:** sync  
**Return:** Book object if found or error message.

### Response

**Example:**

```json
{
    "id": 1,
    "title": "Book Title",
    "author": "Author Name",
    "published_date": "2023-01-15",
    "genre": "Fiction"
}
```

## POST /v1/books

**Description:** Create a new book.  
**Access:** open  
**Operation:** sync  
**Return:** Created book object or error message.  

### Request

**Example:**

```json
{
    "title": "New Book",
    "author": "New Author",
    "published_date": "2023-07-20",
    "genre": "Sci-Fi"
}
```
### Response

**Example:**

```json
{
    "id": 3,
    "title": "New Book",
    "author": "New Author",
    "published_date": "2023-07-20",
    "genre": "Sci-Fi"
}
```

## DELETE /v1/books/{id}

**Description:** Delete a book by its ID.  
**Access:** open  
**Operation:** sync  
**Return:** Success message or error message.

### Response

**Example:**

```json
{
    "message": "Book deleted successfully"
}
```

# Database Description

The database contains information about books, including their title, author, published date, and genre. The database uses MySQL as the backend.

## Table: `book`

The "book" table is the main component of the database, responsible for storing details about various books.

- `id`: This column is of type `INT` and is defined as an auto-incrementing identifier. It is the primary key of the table, ensuring each book entry has a unique identifier.
- `title`: A non-null string column (VARCHAR(128)) holding the title of a book.
- `author`: A non-null string column (VARCHAR(255)) indicating the author of the book.
- `published_date`: A non-null date column (DATE) representing the date of publication.
- `genre`: A non-null string column (VARCHAR(64)) specifying the genre of the book.
