package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// part 2 task 2:Implement the custom response writer

type responseWriter struct {
	http.ResponseWriter     // embed the real writer
	statusCode          int // captured status code
}

// WriteHeader intercepts the status code before forwarding it.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// part 3 task 1: Implement the logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the real writer with our custom one, pre-set to 200
		rw := &responseWriter{w, http.StatusOK}
		// Pass the wrapped writer to the next handler
		next.ServeHTTP(rw, r)
		// The handler has now finished. rw.statusCode holds the real status.
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

//checkin questions answer: each request passes through the middleware once, regardless of which handler responds,
//that is because the middleware wraps the entire router, not individual handlers.

// part 4 task 1: define the application struct
type application struct {
	logger *slog.Logger
}

//checkin question answer: Because grouping dependencies in a struct makes your code more modular, testable and easier to extend.
//handlers dont rely on global state, they access shared resources through the struct, which is more idiomatic Go pattern.

//part 5

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
	app.logger.Info("healthcheck handler called")
}

// List all books
func (app *application) listBooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "list of books (coming soon)\n")
	app.logger.Info("listBooks handler called")
}

// Get a single book by ID
func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "get book with id: %s\n", id)
	app.logger.Info("getBook handler called", "id", id)
}

// Create a new book
func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // <-- must be 201
	fmt.Fprintf(w, "book created (coming soon)\n")
	app.logger.Info("createBook handler called")
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.WriteHeader(http.StatusNoContent)
	app.logger.Info("deleteBook handler called", "id", id)
}

//checkin question answer: HTTP 204 means "no content", so it tells the client the request was successful,
//but there's nothing to return in the reponse body, therefore adding a body would violate the HTTP spec.

//part 6 task 1: compelte main()

func main() {
	// Create a structured logger writing to stdout
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Inject the logger into the application struct
	app := &application{
		logger: logger,
	}
	// Register routes with method-qualified patterns (Go 1.22+)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheck)
	mux.HandleFunc("GET /v1/books", app.listBooks)
	mux.HandleFunc("GET /v1/books/{id}", app.getBook)
	mux.HandleFunc("POST /v1/books", app.createBook)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBook)
	// Log a message before the server starts
	logger.Info("starting server", "addr", ":4000")
	// Wrap the entire router with logging middleware
	err := http.ListenAndServe(":4000", loggingMiddleware(mux))
	log.Fatal(err)
}
