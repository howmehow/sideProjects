package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Books struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Writer *Writer `json:"writer"`
}
type Writer struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Books

func main() {
	r := mux.NewRouter()
	books = append(books, Books{ID: "1", ISBN: "448743", Title: "Book One", Writer: &Writer{Firstname: "Michael", Lastname: "Scott"}})
	books = append(books, Books{ID: "2", ISBN: "12313133", Title: "Book2", Writer: &Writer{Firstname: "Jim", Lastname: "Halpert"}})
	books = append(books, Books{ID: "3", ISBN: "12313131", Title: "Book three", Writer: &Writer{Firstname: "Jimy", Lastname: "Halpert"}})
	books = append(books, Books{ID: "4", ISBN: "1231313123", Title: "Book quattro", Writer: &Writer{Firstname: "Jimmy", Lastname: "Halpert"}})
	books = append(books, Books{ID: "5", ISBN: "12313133123123", Title: "Book funf", Writer: &Writer{Firstname: "Jims", Lastname: "Halpert"}})
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", removeBook).Methods("DELETE")

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all books")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		return
	}
}
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				return
			}
			return
		}
	}

}
func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Books
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000))
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Books
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, item := range books {
		if item.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		return
	}
	return
}
