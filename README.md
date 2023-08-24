
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