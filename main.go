package main

import (
	"log"
	"net/http"
)

// Home handler that writes a byte slice to response body
func home(w http.ResponseWriter, r *http.Request) {
	// only match exactly "/", return 404 otherwise
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST") // has to be done BEFORE `w.WriteHeader` and `w.Write`
		// w.Header()["Date"] = nil                                         // Del() method doesn't remove system-gen headers
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // shortcut for `w.WriteHeader` and `w.Write`
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home) // subtree path like "/**" or "/static/**"
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server at :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
