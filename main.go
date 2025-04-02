package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Book struct {
	BookID          string  `json:"bookId"`
	AuthorID        string  `json:"authorId"`
	PublisherID     string  `json:"publisherId"`
	Title           string  `json:"title"`
	PublicationDate string  `json:"publicationDate"`
	ISBN            string  `json:"isbn"`
	Pages           int     `json:"pages"`
	Genre           string  `json:"genre"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Quantity        int     `json:"quantity"`
}

const filePath = "books.json"

func readBooks() []Book {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return []Book{}
	}
	var books []Book
	json.Unmarshal(file, &books)
	return books
}

func writeBooks(books []Book) {
	file, _ := json.MarshalIndent(books, "", "  ")
	ioutil.WriteFile(filePath, file, 0644)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books := readBooks()
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId := params["id"]
	books := readBooks()
	for _, book := range books {
		if book.BookID == bookId {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	books := readBooks()
	books = append(books, book)
	writeBooks(books)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId := params["id"]
	var updatedBook Book
	_ = json.NewDecoder(r.Body).Decode(&updatedBook)
	books := readBooks()
	for i, book := range books {
		if book.BookID == bookId {
			books[i] = updatedBook
			writeBooks(books)
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId := params["id"]
	books := readBooks()
	for i, book := range books {
		if book.BookID == bookId {
			books = append(books[:i], books[i+1:]...)
			writeBooks(books)
			fmt.Fprintf(w, "Book with ID %s deleted", bookId)
			return
		}
	}
	http.NotFound(w, r)
}

func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	books := readBooks()
	numWorkers := 2 // You can tune this number
	chunkSize := (len(books) + numWorkers - 1) / numWorkers
	resultsChan := make(chan []Book, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(books) {
			end = len(books)
		}
		// Start goroutine for each chunk
		go func(subset []Book) {
			var matches []Book
			for _, book := range subset {
				if strings.Contains(strings.ToLower(book.Title), query) || strings.Contains(strings.ToLower(book.Description), query) {
					matches = append(matches, book)
				}
			}
			resultsChan <- matches
		}(books[start:end])
	}

	// Collect results from all goroutines
	var finalResults []Book
	for i := 0; i < numWorkers; i++ {
		matches := <-resultsChan
		finalResults = append(finalResults, matches...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResults)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	r.HandleFunc("/books/search", searchBooks).Methods("GET")

	log.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", r)
}
