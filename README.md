# REST-API-for-Book-shop (GoLang)
A REST API for managing books using Go and a text file as the data persistence layer.

## Features
CRUD operations for books

Search books by title and description (case-insensitive)

Concurrency for optimized search

Docker containerization

Unit testing for one endpoint

Pagination support

# Getting Started
## Prerequisites
Go 1.21+

Docker (Optional)

Minikube (Optional for Kubernetes deployment)

## Installation
Clone the repository:

git clone <repository-url>
cd bookapi

Install dependencies:

go mod tidy
Running the API
Without Docker

Start the API:
go run main.go
The server will run on http://localhost:8080.

# Method	Endpoint	Description
GET	/books	Get all books (supports pagination)

POST	/books	Create a new book

GET	/books/{id}	Get a book by ID

PUT	/books/{id}	Update a book by ID

DELETE	/books/{id}	Delete a book by ID

GET	/books/search?q=text	Search books by title & description

updated
