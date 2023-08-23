DROP TABLE IF EXISTS book;
CREATE TABLE book (
  id              INT AUTO_INCREMENT NOT NULL,
  title           VARCHAR(128) NOT NULL,
  author          VARCHAR(255) NOT NULL,
  published_date  DATE NOT NULL,
  genre           VARCHAR(64) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO book
  (title, author, published_date, genre)
VALUES
  ('The Great Gatsby', 'F. Scott Fitzgerald', '1925-04-10', 'Classic'),
  ('To Kill a Mockingbird', 'Harper Lee', '1960-07-11', 'Fiction'),
  ('1984', 'George Orwell', '1949-06-08', 'Dystopian'),
  ('Pride and Prejudice', 'Jane Austen', '1813-01-28', 'Romance');

